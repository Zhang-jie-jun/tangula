<template>
  <div>
    <Card title="公共平台列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="12">
              <Row>
                <Col span="4" style="margin-right: 10px;">
                  <Select v-model="searchObj.type" placeholder="选择类型">
                    <Option v-for="(value, key) in typeMapQuery" :value="key" :key="value" >{{ value }}</Option>
                  </Select>
                </Col>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.platName"/>
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
                  @click="platDetail(row)"
                >查看
                </Button>
                <Poptip v-if="access === 'ADMIN'|| access === 'SUPER_ADMIN'"
                  confirm
                  @on-ok="destroy(row)"
                  title="确定要删除平台吗？">
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

    <Modal draggable title="平台信息" width="800px" footer-hide v-model="showDetail">
      <platform-info :b-id=curBuildId :typeMap="typeMap" ref="platInfo"></platform-info>
    </Modal>
  </div>
</template>
<script>
import { httpRequest } from '@/libs/tools.js'
import {
  getPlatformTypes,
  crtPlat,
  getPlatList,
  editPlat,
  publishPlat,
  deletePlat, editHost
} from '@/api/env.js'
import { mapState } from 'vuex'
import platformInfo from "_c/deploy/platform-info";
// import DataFormat from '../../utils/dataFormat.js'

export default {
  name: 'platformPublic',
  // 子组件自己的数据定义，可读可写
  components: {
    platformInfo
  },
  data () {
    return {
      // 搜索
      searchObj: {
        platName: '',
        type:''
      },
      formCrtPlat: {
        type: '',
        name: '',
        ip: '',
        port: '',
        username: '',
        password: ''
      },
      formEditPlat: {
        id: '',
        type: '',
        name: '',
        desc: '',
        ip: '',
        port: '',
        username: '',
        password: ''
      },
      typeMap: {},
      typeMapQuery: {0:'全部'},
      crtModal: false,
      editModal: false,
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        type: [{ required: true, message: '请选择', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入', trigger: 'blur' }],
        name: [{ required: true, message: '请输入', trigger: 'blur' }],
        username: [{ required: true, message: '请输入', trigger: 'blur' }],
        port: [{ required: true, message: '请输入', trigger: 'blur'}],
        password: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      curBuildId: '',
      showDetail: false,
      showSpin: false,
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [
      ],
      columnsTitle: [
        {
          title: 'ID',
          width: 55,
          key: 'id',
          align: 'center'
        },
        {
          title: '名称',
          key: 'name',
          width: 150,
          align: 'center',
          resizable: true,
          render: (h, params) => {
            let arr = []
            arr.push(
              h('span', {
                style: { cursor: 'pointer' }
              },
              [
                h('span',
                  {
                    style: {
                      display: 'inline-block',
                      width: '100%',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap'
                    },
                    domProps: {
                      title: params.row.name
                    }
                  },
                  params.row.name)
              ]
              )

            )
            return h('span', arr)
          }
        },
        {
          title: '类型',
          key: 'type',
          width: 150,
          align: 'center',
          render: (h, params) => {
            return h('span', {
            }, this.typeMap[params.row.type])
          }
        },
        {
          title: 'IP',
          key: 'ip',
          width: 125,
          align: 'center'
        },
        {
          title: '版本',
          key: 'version',
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
          width: 100,
          align: 'center'
        },

        {
          title: '操作',
          key: '',
          width: 220,
          align: 'left',
          slot: 'action'
        }
      ]
    }
  },
  computed: {
    ...mapState({
      access: state => state.user.access,
      userName: state => state.user.userName
    })
  },
  mounted: function () {
    this.getPlatTypes()
    this.getPublicData()
  },
  methods: {
    getPlatTypes () {
      httpRequest(this, getPlatformTypes, []).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        this.typeMap = res.response
        Object.assign(this.typeMapQuery, res.response)
      })
    },

    getPublicData (page = 1) {
      this.showSpin = true
      const params = Object.assign({}, {
        filter: this.searchObj.platName,
        type: this.searchObj.type,
        auth: 222
      })

      httpRequest(this, getPlatList, [page, params]).then(res => {
        console.log('res', res)
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
    openEdit (row) {
      this.editModal = true
      this.formEditPlat.id = row.id
      this.formEditPlat.name = row.name
      this.formEditPlat.ip = row.ip
      this.formEditPlat.port = String(row.port) //此处转string是为了必填校验
      this.formEditPlat.desc = row.desc
      this.formEditPlat.username = row.username
      this.formEditPlat.password = ''
    },

    platDetail (row) {
      this.showDetail = true
      this.curBuildId = row.id
      this.$refs.platInfo.detail(row.id)
    },
    crtPlatform (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            name: this.formCrtPlat.name,
            desc: this.formCrtPlat.desc,
            type: Number(this.formCrtPlat.type),
            ip: this.formCrtPlat.ip,
            port: Number(this.formCrtPlat.port),
            username: this.formCrtPlat.username,
            password: this.formCrtPlat.password,
          }
          crtPlat(param).then(({ data: { code, message } }) => {
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
            this.getPublicData()
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

    editPlat (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            id: this.formEditPlat.id,
            name: this.formEditPlat.name,
            desc: this.formEditPlat.desc,
            type: Number(this.formEditPlat.type),
            ip: this.formEditPlat.ip,
            port: Number(this.formEditPlat.port),
            username: this.formEditPlat.username,
            password: this.formEditPlat.password,
          }
          editPlat(param).then(({ data: { code, message } }) => {
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
            this.editModal = false
            this.getPublicData()
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
    // 发布镜像
    deploy (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, publishPlat, [params]).then(res => {
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
        this.getPublicData()
      })
    },
    destroy (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, deletePlat, [params]).then(res => {
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
        this.getPublicData()
      })
    },
    cancel () {
    },
    // 搜索
    search () {
      this.getPublicData()
      this.currentPage = 1
    },
    // 翻页
    changePage (page) {
      this.getPublicData(page)
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
