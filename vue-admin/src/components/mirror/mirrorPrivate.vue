<template>
  <div>
    <Card title="私有镜像列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="12">
              <Row>
                <Col span="4" style="margin-right: 10px;">
                  <Select v-model="searchObj.dataType" placeholder="选择数据类型">
                    <Option v-for="(value, key) in typeMapQuery" :value="key" :key="value" >{{ value }}</Option>
                  </Select>
                </Col>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.backupName"/>
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
                  @click="detail(row)"
                >查看
                </Button>
                <Button
                  type="warning"
                  size="small"
                  style="margin-right: 5px"
                  @click="openEdit(row)"
                >编辑
                </Button>
                <Button
                  v-if="row.status===1024"
                  type="primary"
                  size="small"
                  style="margin-right: 5px"
                  @click="openBackup(row)"
                >生成副本</Button>
                <Button
                  v-else disabled
                  type="primary"
                  size="small"
                  style="margin-right: 5px"
                >生成副本</Button>
                <Poptip
                  v-if="row.status===1024"
                  confirm
                  title="确定要发布吗？"
                  @on-ok="deployMirror(row)"
                >
                  <Button
                    type="info"
                    size="small"
                    style="margin-right: 5px"
                  >发布
                  </Button>
                </Poptip>
                <Poptip
                  v-else
                  confirm
                  title="确定要发布吗？"
                  @on-ok="deployMirror(row)"
                >
                  <Button
                    type="info" disabled
                    size="small"
                    style="margin-right: 5px"
                  >发布
                  </Button>
                </Poptip>
                <Poptip
                  v-if="row.status===1024"
                  confirm
                  title="确定要销毁吗？"
                  @on-ok="destroyMirror(row)"
                >
                  <Button
                    type="error"
                    size="small"
                    style="margin-right: 5px"
                  >删除
                  </Button>
                </Poptip>
                <Poptip
                  v-else
                  confirm
                  title="确定要销毁吗？"
                  @on-ok="destroyMirror(row)"
                >
                  <Button
                    type="error" disabled
                    size="small"
                    style="margin-right: 5px"
                  >删除
                  </Button>
                </Poptip>
              </div>
            </template>
          </Table>
          <div class="table-page">
            <div class="table-page-position">
              <Page :total="total" :current="currentPage" show-total show-elevator show-sizer placement="top" @on-change="changePage" @on-page-size-change="changeSize" ></Page>
            </div>
          </div>
        </Row>
      </div>
      <Spin fix v-if="showSpin">
        <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
        <div>Loading</div>
      </Spin>
    </Card>
    <Modal
      v-model="backupModal"
      draggable
      title="生成副本"
      @on-cancel="cancel">
      <Form :rules="ruleMirrorValidate" ref="formBackup" :model="formBackup">
        <FormItem label="名称" prop="name">
          <Input v-model="formBackup.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formBackup.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" :disabled="isdisabledFn" v-preventReClick @click="crtBackup('formBackup')">确定</Button>
      </div>
    </Modal>
    <Modal
      v-model="editModal"
      draggable
      title="编辑镜像"
      @on-cancel="cancel">
      <Form  ref="formEditPlat" :model="formEditPlat">
        <FormItem label="名称" prop="desc">
          <Input v-model="formEditPlat.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formEditPlat.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick  @click="editImage('formEditPlat')">确定</Button>
      </div>
    </Modal>
    <Modal draggable title="镜像信息" width="800px" footer-hide v-model="showDetail">
      <mirror-info :b-id=curBuildId :typeMap="typeMap" ref="mirrorInfo"></mirror-info>
    </Modal>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import {
  choiceImageStateColor, choiceTypeColor,
  deleteMirror, editImg,
  getImageStatus,
  getMirrorList,
  getTypes,
  newBackup,
  publishMirror,
  toGB
} from '@/api/backup.js'
import {httpRequest} from "@/libs/tools";
import mirrorInfo from "_c/deploy/mirror-info";

export default {
  name: 'mirrorPrivate',
  // 子组件自己的数据定义，可读可写
  components: {
    mirrorInfo
  },
  data () {
    return {
      isdisabledFn: false, // true 只看进行中的计划
      // 搜索
      searchObj: {
        backupName: '',
        dataType: ''
      },
      formCrtPlat: {
        type: '',
        name: '',
        ip: '',
        port: '',
        userName: '',
        password: ''
      },
      formEditPlat: {
        name: '',
        desc: ''
      },
      editModal: false,
      showDetail: false,
      showSpin: false,
      typeList: ['VMware', 'FusionComputer', 'H3C CAS', 'Hyper-v'],
      curBuildId:'',
      typeMap: {},
      typeMapQuery: {0:'全部'},
      crtModal: false,
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        type: [{ required: true, message: '请选择', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入', trigger: 'blur' }],
        name: [{ required: true, message: '请输入', trigger: 'blur' }],
        userName: [{ required: true, message: '请输入', trigger: 'blur' }],
        password: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      ruleMirrorValidate: {
        name: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      formBackup: {
        mirrorId: '',
        name: '',
        desc: ''
      },
      searchChanged: false,
      backupModal: false,
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      pageSize:10,
      orders: [],
      columnsTitle: [
        {
          title: 'ID',
          key: 'id',
          width: 60,
          align: 'left'
        },
        {
          title: '名称',
          key: 'name',
          align: 'center',
          resizable: true,
          width:200,
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
          title: '大小(GB)',
          key: '',
          width: 90,
          align: 'center',
          render: (h, params) => {
            return h('span', {
            }, toGB(params.row.size))
          }
        },
        {
          title: 'UUID',
          key: '',
          align: 'center',
          render: (h, params) => {
            return h('span', {
            }, params.row.uuid)
          }
        },
        {
          title: '数据类型',
          key: 'type',
          align: 'center',
          render: (h, params) => {
            let color = choiceTypeColor(params.row.type)
            return h('span', {
              style: {
                color: color,
              }
            }, this.typeMap[params.row.type])
          }
        },
        {
          title: '状态',
          width: 130,
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = getImageStatus(params.row.status)
            let color = choiceImageStateColor(params.row.status)
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
          title: '创建时间',
          key: 'created_time',
          width:100,
          align: 'center'
        },

        {
          title: '操作',
          width: 300,
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
    this.getPlatTypes()
    this.getData()
  },
  methods: {
    openCreate () {
      this.crtModal = true
    },
    getPlatTypes () {
      httpRequest(this, getTypes, []).then(res => {
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
    getData (page = 1,pageSize=this.pageSize) {
      this.showSpin = true
      const params = Object.assign({}, {
        auth: 111,
        filter: this.searchObj.backupName,
        type: this.searchObj.dataType
      })

      httpRequest(this, getMirrorList, [page,pageSize, params]).then(res => {
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
    crtBackup (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn=true
          const param = {
            id: this.formBackup.mirrorId,
            name: this.formBackup.name,
            desc: this.formBackup.desc
          }
          this.$Notice.info({
            title: '请求中……',
            desc: '正在生成副本，请等待',
            duration:10
          })
          newBackup(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              this.isdisabledFn=false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '操作成功'
            })
            this.isdisabledFn=false
            this.backupModal = false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
              this.isdisabledFn=false
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    openEdit (row) {
      this.curBuildId = row.id
      this.formEditPlat.desc = row.desc
      this.editModal = true
    },
    editImage (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            id: this.curBuildId,
            name: this.formEditPlat.name,
            desc: this.formEditPlat.desc
          }
          editImg(param).then(({ data: { code, message } }) => {
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
    cancel () {
    },
    openBackup (row) {
      this.backupModal = true
      this.formBackup.mirrorId = row.id
    },
    // 搜索
    search () {
      this.currentPage = 1
      this.getData()
    },
    // 翻页
    changePage (page) {
      this.currentPage = page
      this.getData(page,this.pageSize)
    },
    changeSize (size) {
      this.pageSize=size
      this.getData(1,size);
    },
    closeCreateModal () {
      this.currentTicketId = ''
      this.showCreateSheet = false
      this.modifyCreate = false
    },
    detail (row) {
      this.showDetail = true
      this.curBuildId = row.id
      this.$refs.mirrorInfo.detail(row.id)
    },
    // 发布镜像
    deployMirror (row) {
      console.log('发布镜像', row)
      const params = Object.assign({}, {
        id: row.id,
      })
      httpRequest(this, publishMirror, [params]).then(res => {
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
          this.getData()
        }
      })
    },
    destroyMirror (row) {
      const params = Object.assign({}, {
        id: row.id,
      })
      this.showSpin = true
      httpRequest(this, deleteMirror, [params]).then(res => {
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
        this.showSpin = false
        this.getData()
      })
    }
  }
}
</script>

<style lang="less">
@import "../../view/env/css/ticket.less";
</style>
