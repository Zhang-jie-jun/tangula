<template>
  <div>
    <Card title="用户列表" icon="ios-send" class="card-min-height">
<!--      <Button  type="info" slot="extra" @click.prevent="openCreate" style="margin-bottom: 20px">创建用户</Button>-->
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="12">
              <Row>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="用户名" v-model="searchObj.userName"/>
                </Col>
                <Button class="right-but" type="primary" @click="search">搜索</Button>
              </Row>
            </Col>
            <Col span="3" style="margin-right: 10px;">

            </Col>
          </Row>
          <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small">
            <template slot-scope="{ row }" slot="action">
              <div>
                <Button
                  type="primary"
                  size="small"
                  style="margin-right: 5px"
                  @click="openModifyUser(row)"
                >编辑</Button>
                <Poptip v-if="row.status===1"
                  confirm
                  title="确定要禁用用户吗？"
                  @on-ok="disableUser(row)"
                >
                  <Button
                    type="error"
                    size="small"
                  >禁用
                  </Button>
                </Poptip>
                <Poptip v-if="row.status===-1"
                        confirm
                        title="确定要启用用户吗？"
                        @on-ok="enableUser(row)"
                >
                  <Button
                    type="success"
                    size="small"
                  >启用
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
        </Row>
      </div>
    </Card>

    <Modal
      v-model="modifyUserModal"
      title="编辑用户"
      @on-cancel="cancel">
      <Form ref="formModifyUser" :model="formModifyUser">
        <FormItem label="请选择角色" prop="formModifyUser">
          <Select filterable v-model="formModifyUser.role" >
            <Option v-for="item in roleList" :value="item" :key="item" :label="item" >{{ item }}</Option>
          </Select>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button
          type="primary"
          @click="setRole"
        >确定</Button>
      </div>
    </Modal>
  </div>
</template>
<script>
// import {
//   generateTicket
// } from '@/api/deploy.js'
import { httpRequest } from '@/libs/tools.js'
import { getUserList, getRoleName, disableUser, enableUser, setRole } from '@/api/user.js'
import { mapState } from 'vuex'
// import DataFormat from '../../utils/dataFormat.js'

export default {
  name: 'users',
  // 子组件自己的数据定义，可读可写
  data () {
    return {
      listType: true, // true 只看进行中的计划
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
          width: 70,
          key: 'id',
          align: 'center'
        },
        {
          title: '用户名',
          key: 'account',
          align: 'center'
        },
        {
          title: '邮箱',
          key: 'mail',
          align: 'center'
        },
        {
          title: '角色',
          key: 'role',
          align: 'center',
          render: (h, params) => {
            let roleName
            let color
            if (params.row.role_id === 1) {
              roleName = '普通用户'
            } else {
              roleName = '系统管理员'
              color = '#FF9933'
            }
            return h('span', {
              style: {
                color: color
              }
            }, roleName)
          }
        },
        {
          title: '状态',
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = params.row.status
            let color
            let statusText
            if (lastDeployState === 1) {
              color = 'success'
              statusText = '启用'
            } else {
              color = 'error'
              statusText = '禁用'
            }
            if (!lastDeployState && !params.row.status) {
              return h('span', '')
            }

            if (lastDeployState) {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color
                }
              }, statusText)
            } else {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color
                }
              }, statusText)
            }
          }
        },
        {
          title: '创建时间',
          key: 'created_time',
          align: 'center'
        },
        {
          title: '最后登录',
          key: 'updated_time',
          align: 'center'
          // render: (h, params) => {
          //   return h('span', DataFormat.getDayTime(params.row.updated_at))
          // }
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
      access: state => state.user.access,
      userName: state => state.user.userName
    })
  },
  mounted: function () {
    this.userData()
  },
  methods: {
    changeListType () {
      this.userData()
    },
    userData (page = 1) {
      const params = Object.assign({}, this.searchObj, {
        filter: this.searchObj.userName
      })
      console.log('params', params)
      httpRequest(this, getUserList, [page, params]).then(res => {
        console.log('res', res)
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        this.orders = res.response.data
        this.total = res.response.totalNum
        this.currentPage = page
      })
    },
    disableUser (row) {
      const params = Object.assign({}, row, {
        id: row.id
      })
      httpRequest(this, disableUser, [params]).then(res => {
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
          setTimeout(() => {
            this.userData()
          }, 300)
        }
      })
    },
    enableUser (row) {
      const params = Object.assign({}, row, {
        id: row.id
      })
      httpRequest(this, enableUser, [params]).then(res => {
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
          setTimeout(() => {
            this.userData()
          }, 300)
        }
      })
    },
    setRole () {
      const params = Object.assign({}, {
        id: this.formModifyUser.userId,
        roleName: this.formModifyUser.role
      })
      httpRequest(this, setRole, [params]).then(res => {
        console.log('res', res)
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
        this.modifyUserModal = false
        this.userData()
      })
    },
    // 打开创建发布计划
    openCreate () {
      this.showCreateSheet = true
    },
    cancel () {
    },
    // 打开编辑用户
    openModifyUser (row) {
      this.modifyUserModal = true
      this.formModifyUser.userId = row.id
      this.formModifyUser.role = getRoleName(row.role)
    },
    // 搜索
    search () {
      this.currentPage = 1
      this.userData()
    },
    // 翻页
    changePage (page) {
      this.currentPage = page
      this.userData(page)
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
