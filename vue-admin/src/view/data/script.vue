<template>
  <div>
    <Card title="脚本列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="4">
              <Button type="info"  @click.prevent="openCreate">上传脚本</Button>
            </Col>
            <Col span="12">
              <Row>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.poolName"/>
                </Col>
                <Button class="right-but" type="primary" @click="search">搜索</Button>
              </Row>
            </Col>
          </Row>
          <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small">
            <template slot-scope="{ row }" slot="action">
              <div>
                <Button
                  type="success"
                  size="small"
                  style="margin-right: 5px"
                  @click="openDetails(row)"
                >查看
                </Button>
                <Button
                  type="info"
                  size="small"
                  style="margin-right: 5px"
                  @click="download(row)"
                >下载
                </Button>
                <Poptip
                        confirm
                        @on-ok="destroy(row)"
                        title="确定要删除脚本吗？">
                  <Button
                    type="error"
                    size="small"
                  >删除
                  </Button>
                </Poptip>
              </div>
            </template>
          </Table>
          <div class="table-page">
            <div class="table-page-position">
              <Page :total="total" :current="currentPage" show-total show-elevator @on-change="changePage"></Page>
            </div>
          </div>
          <Spin fix v-if="showSpin">
            <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
            <div>Loading</div>
          </Spin>
        </Row>
      </div>
    </Card>

    <Modal
      v-model="crtModal"
      draggable
      title="上传脚本"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formCrtPlat" :model="formCrtPlat">
        <FormItem label="描述" prop="desc">
          <Input v-model="formCrtPlat.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>

      <Upload ref="upload"
              type="drag"
              style="width: 200px"
              :before-upload="beforeUpload"
              :on-success="successUpload"
              :on-error="failUpload"
              :on-preview="clickFile"
              :data="upParam"
              :action=upUrl
              :headers="upHeaders"
              :max-size="10000"
      >
        <div style="padding: 20px 0">
          <Icon type="ios-cloud-upload" size="22" style="color: #3399ff"></Icon>
          <p>上传脚本(小于10M)</p>
        </div>
      </Upload>
      <div slot="footer">
        <Button  @click="closeModal">关闭</Button>
      </div>
    </Modal>
    <script-details ref="scriptDetailsCom"></script-details>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import {crtPool, getPoolList, deletePool} from '@/api/env.js'
import {httpRequest} from "@/libs/tools";
import poolInfo from "_c/deploy/pool-info";
import ScriptDetails from '_c/deploy/script-details.vue'
import {deleteScript, downloadScript, getScriptList} from "@/api/backup";
import config from "@/config";
import {getToken} from "@/libs/util";
export default {
  name: 'pool',
  // 子组件自己的数据定义，可读可写
  components: {
    ScriptDetails
  },
  data () {
    return {
      listType: true, // true 只看进行中的计划
      scriptName:'',
      upHeaders:{
        Authorization : 'Tangula ' + getToken()
      },
      upUrl:config.baseUrl.pro+'/script/upload',
      upParam:{},
      // 搜索
      searchObj: {
        poolName: ''
      },
      formCrtPlat: {
        type: '',
        name: '',
        desc: ''
      },
      showDetail: false,
      curBuildId: '',
      crtModal: false,
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        name: [{ required: true, message: '请选择', trigger: 'blur' }]
      },
      showSpin: false,

      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [
      ],
      columnsTitle: [
        {
          title: '名称',
          key: 'name',
          align: 'center'
        },
        {
          title: '描述',
          key: 'desc',
          align: 'center'
        },
        {
          title: '创建人',
          key: 'create_user',
          align: 'center'
        },
        {
          title: '创建时间',
          key: 'created_time',
          align: 'center'
        },

        {
          title: '操作',
          key: '',
          align: 'left',
          slot: 'action'
        }
      ]
    }
  },
  computed: {
    ...mapState({
      permits: state => state.user.permits,
      access: state => state.user.access,
      userName: state => state.user.userName
    })
  },
  mounted: function () {
    this.getData()
  },
  methods: {
    getData (page = 1) {
      this.showSpin = true
      const params = Object.assign({}, {
        filter: this.searchObj.poolName,
      })
      httpRequest(this, getScriptList, [page, params]).then(res => {
        this.showSpin = false
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if(res.response.data){
          this.orders = res.response.data
          this.total = res.response.totalNum
        }else {
          this.orders = []
          this.total = 0
        }
        this.currentPage = page

      })
    },
    openCreate () {
      this.crtModal = true
    },
    closeModal(){
      this.crtModal = false
    },
    destroy (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, deleteScript, [params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
        } else {
          this.$Notice.success({
            title: '成功',
            desc: '请求成功'
          })
        }
        this.getData()
      })
    },
    poolDetail (row) {
      this.showDetail = true
      this.curBuildId = row.id
      this.$refs.poolInfo.detail(row.id)
    },
    beforeUpload(file){
      this.scriptName=''
      this.$refs.upload.clearFiles()
      Object.assign(
        this.upParam, {
          desc: this.formCrtPlat.desc
        })
    },
    successUpload(response, file, fileList){
      this.scriptName = file.name;
      this.getData()

    },
    failUpload (err, file, fileList){
      this.$Message.error(err)
    },
    clickFile (file) {
      console.log(this.scriptName)
      //window.open(config.baseUrl.pro+'/tangula/download/script/'+file.name)
    },
    cancel () {
    },
    // 搜索
    search () {
      this.getData()
      this.currentPage = 1
    },
    // 翻页
    changePage (page) {
      this.getData(page)
      this.currentPage = page
    },
    closeCreateModal () {
      this.currentTicketId = ''
      this.showCreateSheet = false
      this.modifyCreate = false
    },
    openDetails (row) {
      this.$refs.scriptDetailsCom.setLogParams(true, row.id)
    },
    download(row){
      window.open(config.baseUrl.pro+'/script/'+row.id+'/download')
    }

  }
}
</script>

<style lang="less">
@import "../../view/env/css/ticket.less";
</style>
