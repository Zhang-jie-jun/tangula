<template>
  <div>
    <Card title="私有主机列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="4">
              <Button type="info"  @click.prevent="openCreate">创建主机</Button>
            </Col>
            <Col span="12">
              <Row>
                <Col span="4" style="margin-right: 10px;">
                  <Select v-model="searchObj.type" placeholder="选择类型">
                    <Option v-for="(value, key) in typeMapQuery" :value="key" :key="value" >{{ value }}</Option>
                  </Select>
                </Col>
                <Col span="6" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.hostName"/>
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
                  @click="hostDetail(row)"
                >查看
                </Button>
                <Button
                  type="primary"
                  size="small"
                  style="margin-right: 5px"
                  @click="openEdit(row)"
                >编辑
                </Button>
                <Poptip
                  confirm
                  title="确定要发布吗？"
                  @on-ok="deploy(row)"
                >
                  <Button
                    type="info"
                    size="small"
                    style="margin-right: 5px"
                  >发布
                  </Button>
                </Poptip>
                <Button v-if="row.status===1"
                  type="info"
                  size="small"
                  style="margin-right: 5px"
                  @click="openDeploy(row)"
                >部署客户端
                </Button>
                <Button v-else
                  type="info" disabled
                  size="small"
                  style="margin-right: 5px"
                >部署客户端
                </Button>
                <Poptip
                  confirm
                  @on-ok="destroy(row)"
                  title="确定要删除主机吗？">
                  <Button
                      type="error"
                      size="small"
                      style="margin-right: 5px"
                  >删除
                  </Button>
                </Poptip>
                <el-button type="success"
                           size="mini"
                           @click="update(row)"
                           circle>
                  刷新
                </el-button>
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
      title="创建主机"
      @on-cancel="cancel">
      <Tooltip  placement="right-start" max-width="1024">
        <Tag color="warning">前置要求</Tag>
        <div slot="content">
          <p>鉴于添加平台的过程中需要对远程主机进行认证操作，请按照以下要求检查目标主机环境:</p>
          <br>
          <p>Linux主机：请确保需要添加的目标主机系统已经安装并启动ssh服务与nfs客户端服务。</p>
          <p>Windows主机：请确保需要添加的目标主机系统已经启动winrm服务与nfs客户端服务。</p>
          <br>
          <p>更多详细说明请参阅使用手册</p>
        </div>
      </Tooltip>
      <Form :rules="ruleValidate" ref="formCrtPlat" :model="formCrtPlat">
        <FormItem label="名称" prop="name">
          <Input v-model="formCrtPlat.name" placeholder="请输入名称" />
        </FormItem>
        <FormItem label="请选择主机类型" prop="type">
          <Select filterable v-model="formCrtPlat.type" >
            <Option v-for="(value, key) in typeMap" :value="key" :key="value" >{{ value }}</Option>
          </Select>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formCrtPlat.desc" placeholder="请输入描述"/>
        </FormItem>
        <FormItem label="IP" prop="ip">
          <Input v-model="formCrtPlat.ip" placeholder="请输入ip"/>
        </FormItem>
        <FormItem label="端口号" prop="port">
          <Input v-model="formCrtPlat.port" placeholder="Linux默认端口22，Windows默认端口5985"/>
        </FormItem>
        <FormItem label="用户名" prop="username">
          <Input v-model="formCrtPlat.username" placeholder="请输入用户名"/>
        </FormItem>
        <FormItem label="密码" prop="password">
          <Input v-model="formCrtPlat.password" type="password" placeholder="请输入密码"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="crtPlat('formCrtPlat')">确定</Button>
      </div>
    </Modal>

    <Modal
      v-model="editModal"
      draggable
      title="编辑主机"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formEditPlat" :model="formEditPlat">
        <FormItem label="名称" prop="name">
          <Input v-model="formEditPlat.name" placeholder="请输入名称" disabled/>
        </FormItem>
        <FormItem label="请选择主机类型" prop="type">
          <Select filterable v-model="formEditPlat.type" >
            <Option v-for="(value, key) in typeMap" :value="key" :key="value" >{{ value }}</Option>
          </Select>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formEditPlat.desc" placeholder="请输入描述"/>
        </FormItem>
        <FormItem label="IP" prop="ip">
          <Input v-model="formEditPlat.ip" placeholder="请输入ip"/>
        </FormItem>
        <FormItem label="端口号" prop="port">
          <Input v-model="formEditPlat.port" placeholder="Linux默认端口22，Windows默认端口5985"/>
        </FormItem>
        <FormItem label="用户名" prop="username">
          <Input v-model="formEditPlat.username" placeholder="请输入用户名"/>
        </FormItem>
        <FormItem label="密码" prop="password">
          <Input v-model="formEditPlat.password" type="password" placeholder="请输入密码"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="editPlat('formEditPlat')">确定</Button>
      </div>
    </Modal>

    <el-dialog
      :visible.sync="deployModal"
      :close-on-click-modal=false
      width="800px"
      title="部署客户端"
      @on-cancel="cancel">
      <el-form :rules="ruleValidate" ref="formDeploy" :model="formDeploy">
        <el-form-item label="安装目录" :label-width="formLabelWidth" prop="baseDir">
          <el-input v-model="formDeploy.baseDir" autocomplete="off" placeholder="请输入安装目录"></el-input>
        </el-form-item>
        <el-form-item label="控制台ip" :label-width="formLabelWidth" prop="serverIp">
          <el-input v-model="formDeploy.serverIp" autocomplete="off" placeholder="请输入ip"></el-input>
        </el-form-item>
        <el-form-item label="客户端应用" :label-width="formLabelWidth" prop="apps">
          <el-select v-model="formDeploy.apps" multiple  placeholder="支持多选">
            <el-option
              v-for="item in appList"
              :key="item"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="安装类型" :label-width="formLabelWidth" prop="abType">
            <el-select v-model="formDeploy.abType"   placeholder="请选择安装类型">
              <el-option
                v-for="item in abTypeList"
                :key="item"
                :label="item"
                :value="item">
              </el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="Basic ftp路径" :label-width="formLabelWidth" prop="basicFtpPath">
          <el-input v-model="formDeploy.basicFtpPath" autocomplete="off" placeholder="ftp路径若未输入则取默认的包"></el-input>
        </el-form-item>
        <el-form-item label="应用ftp路径" :label-width="formLabelWidth" prop="appFtpPath">
          <el-input v-model="formDeploy.appFtpPath" autocomplete="off" placeholder="ftp路径若未输入则取默认的包;多个路径用分号分隔"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="installClient('formDeploy')">确定</Button>
      </div>
    </el-dialog>

    <Modal draggable title="主机信息" width="1024px" footer-hide v-model="showDetail">
      <build-info :b-id=curBuildId :typeMap="typeMap" ref="buildInfo"></build-info>
    </Modal>
  </div>
</template>
<script>

import { httpRequest } from '@/libs/tools.js'
import { mapState } from 'vuex'
import {
  getHostList,
  getHostTypes,
  crtHost,
  editHost,
  publishHost,
  deletehHost,
  updateHost,
  deployClient,
  getDeployStatus, choiceDeployStateColor
} from '@/api/env.js'
import buildInfo from '_c/deploy/build-info'
export default {
  name: 'machinePrivate',
  // 子组件自己的数据定义，可读可写
  components: {
    buildInfo
  },
  data () {
    return {
      // 搜索
      searchObj: {
        hostName: '',
        type: ''
      },
      showSpin: false,
      formCrtPlat: {
        type: '',
        name: '',
        desc: '',
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
      formDeploy: {
        id: '',
        serverIp: '',
        apps: '',
        baseDir: '',
        abType: '',
        basicFtpPath: '',
        appFtpPath: ''
      },
      typeMap: {},
      typeMapQuery: { 0: '全部' },
      crtModal: false,
      editModal: false,
      deployModal: false,
      formLabelWidth: '120px',
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        type: [{ required: true, message: '请选择', trigger: 'blur' }],
        ip: [{ required: true, message: '请输入', trigger: 'blur' }],
        name: [{ required: true, message: '请输入', trigger: 'blur' }],
        port: [{ required: true, message: '请输入', trigger: 'blur' }],
        username: [{ required: true, message: '请输入', trigger: 'blur' }],
        password: [{ required: true, message: '请输入', trigger: 'blur' }],
        serverIp: [{ required: true, message: '请输入', trigger: 'blur' }],
        apps: [{ required: true, message: '请输入', trigger: 'blur' }],
        baseDir: [{ required: true, message: '请输入', trigger: 'blur' }],
        abType: [{ required: true, message: '请选择', trigger: 'blur' }]
      },
      showDetail: false,
      curBuildId: '',
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      orders: [],
      appList: [
        'VMware',
        'FusionCompute',
        'ZStack',
        'XenServer',
        'SmartX',
        'H3CCAS',
        'QingCloud'
      ],
      abTypeList: [
        'AB'
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
          align: 'center',
          width: 150,
          resizable: true,
          render: (h, params) => {
            let arr = []
            arr.push(
              h('span', { style: { cursor: 'pointer' } },
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
          width: 100,
          key: 'type',
          align: 'center',
          render: (h, params) => {
            return h('span', {
            }, this.typeMap[params.row.type])
          }
        },
        {
          title: 'IP',
          width: 160,
          key: 'ip',
          align: 'left',
          render: (h, params) => {
            let atts = []
            let color
            let text
            if (params.row.status === 1) {
              color = '#66CC00'
              text = '已连通'
            } else {
              color = '#CCCCCC'
              text = '未连通'
            }
            atts.push(
              h('Tooltip', {
                props: { content: text, placement: 'top' },
                style: { cursor: 'pointer', color: color, marginRight: 15 } }, [h('Icon', { props: { type: 'md-desktop' }, style: { fontsize: '24px' } })]
              )
            )
            atts.push(params.row.ip)
            return h('span', atts)
          }
        },
        {
          title: '系统',
          key: 'os',
          align: 'center'
        },
        {
          title: '客户端部署状态',
          width: 130,
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = getDeployStatus(params.row.deployStatus)
            let color = choiceDeployStateColor(params.row.deployStatus)
            if (!lastDeployState && !params.row.deployStatus) {
              return h('span', '')
            }

            if (lastDeployState) {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color,
                  fade: true
                }
              }, lastDeployState)
            } else {
              return h('Tag', {
                props: {
                  type: 'dot',
                  color: color,
                  fade: true
                }
              }, params.row.deployStatus)
            }
          }
        },
        {
          title: '创建人',
          key: 'create_user',
          align: 'center'
        },
        {
          title: '创建时间',
          width: 100,
          key: 'created_time',
          align: 'center'
        },
        {
          title: '更新时间',
          width: 100,
          key: 'updated_time',
          align: 'center'
        },
        {
          title: '操作',
          key: '',
          width: 350,
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
    this.getHostTypes()
    this.getPrivateData()
  },
  methods: {
    getPrivateData (page = 1) {
      this.showSpin = true
      const params = Object.assign({}, {
        filter: this.searchObj.hostName,
        type: this.searchObj.type,
        auth: 111
      })
      console.log('params', params)
      httpRequest(this, getHostList, [page, params]).then(res => {
        this.showSpin = false
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if (res.response.data) {
          this.orders = res.response.data
          this.total = res.response.totalNum
        } else {
          this.orders = []
          this.total = 0
        }
        this.currentPage = page
      })
    },
    getHostTypes () {
      httpRequest(this, getHostTypes, []).then(res => {
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
    openCreate () {
      this.crtModal = true
      this.formCrtPlat.password = ''
    },
    openEdit (row) {
      this.editModal = true
      this.formEditPlat.id = row.id
      this.formEditPlat.name = row.name
      this.formEditPlat.ip = row.ip
      this.formEditPlat.port = String(row.port) // 此处转string是为了必填校验
      this.formEditPlat.desc = row.desc
      this.formEditPlat.username = row.username
      this.formEditPlat.password = ''
    },
    openDeploy (row) {
      this.deployModal = true
      this.formDeploy.id = row.id
    },
    crtPlat (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            name: this.formCrtPlat.name,
            desc: this.formCrtPlat.desc,
            type: Number(this.formCrtPlat.type),
            ip: this.formCrtPlat.ip,
            port: Number(this.formCrtPlat.port),
            username: this.formCrtPlat.username,
            password: this.formCrtPlat.password
          }
          crtHost(param).then(({ data: { code, message } }) => {
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
            this.getPrivateData()
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
            password: this.formEditPlat.password
          }
          editHost(param).then(({ data: { code, message } }) => {
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
            this.getPrivateData()
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
    installClient (name) {
      this.$refs[name].validate((valid) => {
        if (!valid) {
          this.$Message.error('请输入必填项!')
        } else {
          let apps_str = this.formDeploy.apps.join(',')
          const param = {
            hostId: this.formDeploy.id,
            baseDir: this.formDeploy.baseDir,
            serverIp: this.formDeploy.serverIp,
            apps: apps_str,
            type: this.formDeploy.abType,
            basicFtpPath: this.formDeploy.basicFtpPath,
            appFtpPath: this.formDeploy.appFtpPath
          }
          console.log(param)
          deployClient(param).then(({ data: { code, message } }) => {
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
            this.deployModal = false
            this.getPrivateData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
            }
          })
        }
      })
    },
    cancel () {
    },
    hostDetail (row) {
      this.showDetail = true
      this.curBuildId = row.id
      this.$refs.buildInfo.detail(row.id)
      this.$refs.buildInfo.getRecordData(1, row.id)
    },
    // 搜索
    search () {
      this.getPrivateData()
      this.currentPage = 1
    },
    // 翻页
    changePage (page) {
      this.getPrivateData(page)
      this.currentPage = page
    },
    closeCreateModal () {
      this.currentTicketId = ''
      this.showCreateSheet = false
      this.modifyCreate = false
    },

    // 更新状态
    update (row) {
      const params = Object.assign({}, {
        id: row.id
      })
      this.$Notice.info({
        title: '请求中……',
        desc: '正在查询主机状态,请稍候'
      })
      httpRequest(this, updateHost, [params]).then(res => {
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
        this.getPrivateData()
      })
    },
    // 发布镜像
    deploy (row) {
      const params = Object.assign({}, {
        id: row.id
      })
      httpRequest(this, publishHost, [params]).then(res => {
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
        this.getPrivateData()
      })
    },
    destroy (row) {
      const params = Object.assign({}, {
        id: row.id
      })
      httpRequest(this, deletehHost, [params]).then(res => {
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
        this.getPrivateData()
      })
    }
  }
}
</script>

<style lang="less">
@import "../../view/env/css/ticket.less";
</style>
