<template>
  <div>
    <Card title="存储池列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="4" v-if="this.access === 'ADMIN'|| this.access === 'SUPER_ADMIN'">
              <Button type="info"  @click.prevent="openCreate">创建池</Button>
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
                  @click="poolDetail(row)"
                >查看
                </Button>
                <Poptip v-if="access === 'ADMIN'|| access === 'SUPER_ADMIN'"
                  confirm
                  @on-ok="destroy(row)"
                  title="确定要删除存储池吗？">
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
      title="创建池"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formCrtPlat" :model="formCrtPlat">
        <FormItem label="名称" prop="name">
          <Input v-model="formCrtPlat.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formCrtPlat.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="crtPlat('formCrtPlat')">确定</Button>
      </div>
    </Modal>
    <Modal draggable title="存储池信息" width="800px" footer-hide v-model="showDetail">
      <pool-info :b-id=curBuildId ref="poolInfo"></pool-info>
    </Modal>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import {crtPool, getPoolList, deletePool} from '@/api/env.js'
import {httpRequest} from "@/libs/tools";
import poolInfo from "_c/deploy/pool-info";
export default {
  name: 'pool',
  // 子组件自己的数据定义，可读可写
  components: {
    poolInfo
  },
  data () {
    return {
      listType: true, // true 只看进行中的计划
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
      httpRequest(this, getPoolList, [page, params]).then(res => {
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
    crtPlat (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            name: this.formCrtPlat.name,
            desc: this.formCrtPlat.desc
          }
          crtPool(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '操作成功'
            })
            this.crtModal = false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },

    destroy (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, deletePool, [params]).then(res => {
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
    }
  }
}
</script>

<style lang="less">
@import "../../view/env/css/ticket.less";
</style>
