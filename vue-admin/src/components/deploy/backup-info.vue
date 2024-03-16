<template>
  <div>
    <Collapse v-model="showPanel">
      <Panel name="record">
        <strong style="color: #3399CC;font-size: 16px">挂载/卸载执行记录</strong>
        <div slot="content">
          <Row>
            <Table searchable border :columns="columnsTitle" :data="orders" stripe size="small"></Table>
            <div class="table-page">
              <div class="table-page-position">
                <Page :total="total" :current="currentPage" show-total show-elevator @on-change="changePage"></Page>
              </div>
            </div>
          </Row>
        </div>
      </Panel>
      <Panel name="record">
        <strong style="color: #3399CC;font-size: 16px">副本信息</strong>
        <p slot="content">
        <Form :label-width="90" label-position="left" class="ivu-form-item-label">
          <Row>
            <Col span="12">
              <FormItem label="副本名">
                <span>{{ detailData.name }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="描述">
                <span>{{ detailData.desc }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="副本uuid">
                <span>{{ detailData.uuid }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="大小">
                <span>{{ detailData.size }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="类型">
                <span>{{ typeText }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="状态">
                <span>{{ statusText }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="挂载对象">
                <span>{{ mountName }}</span>
              </FormItem>
            </Col>
<!--            <Col span="12">-->
<!--              <FormItem label="挂载点">-->
<!--                <span>{{ mountPoint }}</span>-->
<!--              </FormItem>-->
<!--            </Col>-->
            <Col span="12">
              <FormItem label="导出路径">
                <span>{{ detailData.export }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="存储池">
                <span>{{ poolData.name }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="存储池uuid">
                <span>{{ poolData.uuid }}</span>
              </FormItem>
            </Col>
          </Row>
          <Row>
            <Col span="12">
              <FormItem label="创建人">
                <span>{{ detailData.create_user }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="创建时间">
                <span>{{ detailData.created_time }}</span>
              </FormItem>
            </Col>
          </Row>
        <Row v-if="(this.type===1005 || this.type===1004) && this.status===4096">
          <Col span="12">
            <FormItem label="配置文件">
              <span v-if="this.isUpd" style="color: #339933">已上传</span>
              <span v-else style="color: #FF0033">未上传</span>
              <Button v-if="this.isUpd"
                size="small"
                type="success"
                style="margin-left: 5px"
                @click="openDetails"
              >查看
              </Button>
            </FormItem>
          </Col>
          <Col span="12">
            <FormItem label="操作">
              <Button type="info"  @click.prevent="openCreate">上传配置文件</Button>
<!--              <Button type="success" style="margin-left: 10px" @click="openDoCas">执行兼容性测试</Button>-->
            </FormItem>
          </Col>
        </Row>
        </Form>
        </p>
      </Panel>
    </Collapse>
    <deploy-log ref="deployLogCom"></deploy-log>

    <Modal
      v-model="crtModal"
      draggable
      title="上传配置文件"
      @on-cancel="cancel">
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
          <p>上传配置文件(小于10M)</p>
        </div>
      </Upload>
      <div slot="footer">
        <Button  @click="closeModal">关闭</Button>
      </div>
    </Modal>

    <Modal
      v-model="doCasModal"
      draggable
      title="执行兼容性测试AT"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formDoCas" :model="formDoCas">
        <FormItem label="控制台ip" prop="ip">
          <Input v-model="formDoCas.ip" placeholder="请输入控制台ip"/>
        </FormItem>
        <FormItem label="测试用例tag" prop="tag">
          <Input v-model="formDoCas.tag" placeholder="测试用例tag"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="doCas('formDoCas')">确定</Button>
      </div>
    </Modal>
    <cas-config-details ref="configDetailsCom"></cas-config-details>
  </div>

</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import {
  backupDetails, doCasJenkins,
  getImageStatus,
  getInstanceList

} from "@/api/backup.js"
import {crtPlat, getHostTypes, getPlatformTypes} from "@/api/env"
import DeployLog from '_c/deploy/show-log.vue'
import {getToken} from "@/libs/util"
import config from "@/config"
import casConfigDetails from '_c/deploy/cas-config-details.vue'
export default {
  name: 'backupInfo',
  components: {
    DeployLog,
    casConfigDetails
  },
  props: {
    bId: {
      type: [Number, String],
      required: false
    },
    typeMap: Object,
    backupId: Number
  },
  data () {
    return {
      curReplicaId:'',
      detailData: {},
      poolData: {},
      platTypeMap: {},
      hostTypeMap: {},
      status:'',
      statusText:'',
      type:'',
      typeText: '',
      targetText: '',
      showDetail: false,
      crtModal: false,
      isUpd:false,
      scriptName:'',
      upHeaders:{
        Authorization : 'Tangula ' + getToken()
      },
      upUrl:config.baseUrl.pro+'/store_pool/replica/uploadFile',
      upParam:{},
      ruleValidate: {
        ip: [{ required: true, message: '请输入', trigger: 'blur' }],
        tag: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      formDoCas: {
        ip: '',
        tag: ''
      },
      doCasModal: false,

      showPanel: 'record',
      mountPoint: '',
      mountName: '',
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
          title: '操作类型',
          key: 'type',
          align: 'center',
          render: (h, params) => {
            let color = this.getInstanceColor(params.row.type)
            return h('span', {
              style: {
                color: color
              }
            }, this.getInstanceType(params.row.type))
          }
        },
        {
          title: '操作平台',
          key: 'target_type',
          align: 'center',
          render: (h, params) => {
            let typeText
            if (params.row.target_type < 50) { // 平台
              typeText = this.platTypeMap[params.row.target_type]
            } else {
              typeText = this.hostTypeMap[params.row.target_type]
            }
            return h('span', {
            }, typeText)
          }
        },

        {
          title: '挂载对象',
          key: 'type',
          align: 'center',
          render: (h, params) => {
            return h('span', {
            }, params.row.targetInfo.name)
          }
        },
        {
          title: '挂载点',
          key: 'mount_point',
          align: 'center'
        },
        {
          title: '状态',
          width: 150,
          key: 'status',
          align: 'center',
          render: (h, params) => {
            let lastDeployState = this.getInstanceStatus(params.row.status)
            let color = this.getInstanceStateColor(params.row.status)
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
          title: '执行时间',
          key: 'created_time',
          align: 'center'
        },
        {
          title: '操作',
          key: '',
          width: 150,
          align: 'center',
          style: 'padding: 0px',
          render: (h, params) => {
            let atts = []
            atts.push(
              h('Button', {
                props: {
                  type: 'success',
                  ghost: ''
                },
                on: {
                  click: () => {
                    this.openLog(params.row)
                  }
                }
              }, '执行日志')
            )

            return h('div', atts)
          }
        }
      ]
    }
  },
  computed: {
    deployColor () {
      return '#9933FF'
    }
  },
  mounted: function () {
    this.getHostTypes()
    this.getPlatTypes()
  },
  methods: {
    detail(_id) {
      this.curReplicaId=_id
      this.showDetail = true
      httpRequest(this, backupDetails, [_id]).then(res => {
        this.detailData = res.response
        this.poolData = res.response.pool
        this.status=res.response.status
        this.statusText = getImageStatus(res.response.status)
        this.type=res.response.type
        this.typeText = this.typeMap[res.response.type]
        this.mountPoint=res.response.mount_info?res.response.mount_info.mount_param:''
        this.mountName=res.response.mount_info?res.response.mount_info.targetInfo.name:''
        this.isUpd=res.response.isUpload
      })
    },
    getRecordData (page = 1,_id) {
      let backId = _id ? _id :this.bId
      const params = Object.assign({}, {
        id: backId
      })

      httpRequest(this, getInstanceList, [page, params]).then(res => {
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

    openDoCas () {
      this.doCasModal = true
    },

    doCas(name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            id: this.curReplicaId,
            consoleIp:this.formDoCas.ip,
            tag:this.formDoCas.tag
          }
          doCasJenkins(param).then(({ data: { code, message,response } }) => {
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
            this.doCasModal = false
            let buildId=Number(response.buildId)+1
            window.open("http://192.168.212.140:8080/jenkins/view/兼容性测试/job/CAS_兼容性测试")

          }).catch(err => {
            if (err && err.response) {
              this.$Notice.error({
                title: '错误',
                desc: err.response
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

    changePage (page) {
      this.currentPage = page
      this.getRecordData(page);
    },
    getInstanceType (val){
      const statusNum = Number(val)
      const statusMap = {
        1: '挂载',
        2: '卸载'
      }
      if (statusMap.hasOwnProperty(statusNum)) {
        return statusMap[statusNum]
      } else {
        return '挂载'
      }
    },

    getInstanceColor (val) {
      const statusNum = Number(val)
      const statusMap = {
        1: '#339933',
        2: '#FF6666'
      }
      if (statusMap.hasOwnProperty(statusNum)) {
        return statusMap[statusNum]
      } else {
        return '#0099CC'
      }
    },
    getPlatTypes () {
      httpRequest(this, getPlatformTypes, []).then(res => {
        if (res.code !== 200) {
          return
        }
        this.platTypeMap = res.response
      })
    },
    getHostTypes () {
      httpRequest(this, getHostTypes, []).then(res => {
        if (res.code !== 200) {
          return
        }
        this.hostTypeMap = res.response
      })
    },
    getInstanceStatus(val){
      const statusNum = Number(val)
      const statusMap = {
        1: 	 '未启动',
        2:	 '准备中',
        4: 	'挂载中',
        8 :	 '挂载成功',
        16:	 '挂载失败',
        32:	 '卸载中',
        64:	 '卸载成功',
        128: '卸载失败',
        256: '异常'
      }
      if (statusMap.hasOwnProperty(statusNum)) {
        return statusMap[statusNum]
      } else {
        return '未启动'
      }
    },
    getInstanceStateColor (state) {
      let color = 'warning'
      switch (state) {
        case 1:
          color = 'warning'
          break
        case 2:
          color = 'primary'
          break
        case 4:
          color = 'primary'
          break
        case 8:
          color = 'success'
          break
        case 16:
          color = 'error'
          break
        case 32:
          color = 'primary'
          break
        case 64:
          color = 'success'
          break
        case 128:
          color = 'error'
          break
        case 256:
          color = 'error'
          break
        default:
          color = 'error'
      }
      return color
    },

    openCreate () {
      this.crtModal = true
    },
    closeModal(){
      this.crtModal = false
    },

    beforeUpload(file){
      this.scriptName=''
      this.$refs.upload.clearFiles()
      Object.assign(
        this.upParam, {
          id: this.curReplicaId
        })
    },
    successUpload(response, file, fileList){
      if (response.code!==200){
        //this.$Message.error('上传失败:'+response.message)
        this.$Message.error({
          content: '上传失败:'+response.message,
          duration: 3
        });
      }else {
        this.$Message.success('上传成功！')
      }
      this.scriptName = file.name;
    },
    failUpload(err, file, fileList){
      this.$Message.error('上传失败！'+err)
    },
    clickFile(file){
      console.log(this.scriptName)
      // window.open(config.baseUrl.pro+'/tangula/download/script/'+file.name)
    },
    openLog (row) {
      this.$refs.deployLogCom.setLogParams(true, row.id)
    },
    openDetails () {
      this.$refs.configDetailsCom.setFileParams(true, this.curReplicaId)
    },
  }
}
</script>
