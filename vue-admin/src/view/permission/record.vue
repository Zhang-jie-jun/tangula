<template>
  <div>
    <Card title="操作记录" icon="ios-send" class="card-min-height">
      <!--      <Button  type="info" slot="extra" @click.prevent="openCreate" style="margin-bottom: 20px">创建用户</Button>-->
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="12">
              <Row>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="操作人" v-model="searchObj.userName"/>
                </Col>
                <Button class="right-but" type="primary" @click="search">搜索</Button>
              </Row>
            </Col>
            <Col span="3" style="margin-right: 10px;">

            </Col>
          </Row>
          <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small">
          </Table>
          <div class="table-page">
            <div class="table-page-position">
              <Page :total="total" :current="currentPage" show-total show-elevator @on-change="changePage"></Page>
            </div>
          </div>
        </Row>
      </div>
      <Spin fix v-if="showSpin">
        <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
        <div>Loading</div>
      </Spin>
    </Card>

    <deploy-log ref="deployLogCom"></deploy-log>
  </div>
</template>
<script>
// import {
//   generateTicket
// } from '@/api/deploy.js'
import { mapState } from 'vuex'
import DeployLog from '_c/deploy/show-log.vue'
import {httpRequest} from "@/libs/tools";
import {
  choiceImageStateColor,
  choiceRecordStateColor,
  getImageStatus,
  getMirrorList,
  getRecordList,
  getRecordStatus
} from "@/api/backup";

export default {
  name: 'record',
  components: {
    DeployLog
  },
  // 子组件自己的数据定义，可读可写
  data () {
    return {
      // 搜索
      searchObj: {
        userName: ''
      },
      formModifyUser: {
        userId: 0,
        role: ''
      },
      roleList: ['系统管理员', '普通用户'],
      modifyUserModal: false,
      usersRole: '',
      searchChanged: false,
      showSpin: false,
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [],
      selIssues: [],
      sqlScripts: '',
      curOrderId: '',
      ticketData: {},
      showTicket: false,
      showCreateTicket: false,
      req: false,
      showCreated: false,
      defaultRows: 6,
      columnsTitle: [
        {
          title: 'ID',
          width: 75,
          key: 'id',
          align: 'center'
        },
        {
          title: '操作',
          key: 'operation',
          align: 'center'
        },
        {
          title: '名称',
          key: 'object',
          align: 'center'
        },
        {
          title: '状态',
          width: 130,
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = getRecordStatus(params.row.status)
            let color = choiceRecordStateColor(params.row.status)
            if (!lastDeployState && !params.row.status) {
              return h('span', '')
            }

            if (lastDeployState) {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color
                }
              }, lastDeployState)
            } else {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color
                }
              }, params.row.status)
            }
          }
        },
        {
          title: '详情',
          key: 'detail',
          align: 'center'
        },
        {
          title: '操作人',
          key: 'user',
          align: 'center'
        },
        {
          title: '操作时间',
          key: 'created_time',
          align: 'center'
        }
        // {
        //   title: '执行日志',
        //   key: '',
        //   align: 'center',
        //   style: 'padding: 0px',
        //   render: (h, params) => {
        //     let atts = []
        //       atts.push(
        //         h('Button', {
        //           props: {
        //             type: 'success',
        //             ghost: ''
        //           },
        //           on: {
        //             click: () => {
        //               this.openLog(params.row)
        //             }
        //           }
        //         }, '查看')
        //       )
        //
        //     return h('div', atts)
        //   }
        // }
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
    this.getData()
  },
  methods: {
    getData(page = 1) {
      this.showSpin = true
      const params = Object.assign({}, {
        user: this.searchObj.userName
      })

      httpRequest(this, getRecordList, [page, params]).then(res => {
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
      this.showCreateSheet = true
    },
    cancel () {
    },
    // 搜索
    search () {
      this.currentPage = 1
      this.getData()
    },
    // 翻页
    changePage (page) {
      this.currentPage = page
      this.getData(page)
    },
    openLog (row) {
      this.$refs.deployLogCom.setLogParams(true, row.id)
    }
  }
}
</script>

<style lang="less">
</style>
