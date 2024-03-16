<template>
  <div>
    <Card title="角色列表" icon="ios-send" class="card-min-height">
      <!--      <Button  type="info" slot="extra" @click.prevent="openCreate" style="margin-bottom: 20px">创建用户</Button>-->
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="12">
              <Row>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="角色名" v-model="searchObj.roleName"/>
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
  </div>
</template>
<script>
// import {
//   generateTicket
// } from '@/api/deploy.js'
import { httpRequest } from '@/libs/tools.js'
import { getRoleList, getRoleName } from '@/api/user.js'
import { mapState } from 'vuex'
// import DataFormat from '../../utils/dataFormat.js'

export default {
  name: 'roles',
  // 子组件自己的数据定义，可读可写
  data () {
    return {
      listType: true, // true 只看进行中的计划
      // 搜索
      searchObj: {
        roleName: ''
      },
      formModifyUser: {
        role: ''
      },
      roleList: [],
      modifyUserModal: false,
      usersRole: '',
      searchChanged: false,

      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [],
      columnsTitle: [
        {
          title: 'ID',
          width: 70,
          key: 'id',
          align: 'center'
        },
        {
          title: '角色名',
          key: 'name',
          align: 'center'
        },
        {
          title: '创建时间',
          key: 'created_time',
          align: 'center'
        }
        // {
        //   title: '操作',
        //   key: '',
        //   align: 'left',
        //   slot: 'action'
        // }
      ]
    }
  },
  computed: {
    ...mapState({
      permits: state => state.user.permits,
      userName: state => state.user.userName
    })
  },
  mounted: function () {
    this.searchObj.biz = this.defBusi === '无' ? '' : this.defBusi
    this.userData()
  },
  methods: {
    changeListType () {
      this.userData()
    },
    userData (page = 1) {
      const params = Object.assign({}, this.searchObj, {
        filter: this.searchObj.roleName
      })
      httpRequest(this, getRoleList, [page, params]).then(res => {
        console.log('res', res)
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        // this.orders = res.response.data
        this.total = res.response.totalNum
        for(let i=0;i<res.response.data.length;i++){
          if(res.response.data[i].name !== '超级管理员'){
            this.orders.push(res.response.data[i])
          }else {
            this.total -=1
          }
        }
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
