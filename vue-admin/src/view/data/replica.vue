<template>
  <div>
    <Card title="副本列表" icon="ios-send" class="card-min-height">
      <!-- 主体部分 -->
      <div>
        <Row>
          <Row class="search-margin-botton">
            <Col span="8">
              <Button type="primary"  @click.prevent="openCreate" style="margin-right: 5px">创建副本</Button>
              <Button type="info"   @click="openMountBatch" style="margin-right: 5px">批量挂载</Button>
              <Button type="warning"  @click="batchUnMount" style="margin-right: 5px">批量卸载</Button>
              <Button type="error"  @click="batchDelete">批量删除</Button>
            </Col>
            <Col span="12">
              <Row>
                <Col span="3" style="margin-right: 10px;">
                  <Select v-model="searchObj.status" placeholder="状态" >
                    <Option
                      v-for="(item, key) in backupStatus"
                      :key="key"
                      :label="item.key"
                      :value="item.value"></Option>
                  </Select>
                </Col>
                <Col span="4" style="margin-right: 10px;">
                  <Select v-model="searchObj.dataType" placeholder="数据类型">
                    <Option v-for="(value, key) in typeMapQuery" :value="key" :key="value" >{{ value }}</Option>
                  </Select>
                </Col>
                <Col span="4" style="margin-right: 10px;">
                  <Input placeholder="名称" v-model="searchObj.backupName"/>
                </Col>
                <Col span="4" style="margin-right: 10px;">
                  <Input placeholder="uuid" v-model="searchObj.uuid"/>
                </Col>
                <Button class="right-but" type="primary" @click="search">搜索</Button>
              </Row>
            </Col>
            <Button type="success" @click="refresh" size="small" style="float:right">刷新</Button>
          </Row>
          <Alert>
            <strong>注意：</strong> Custom Image类型副本挂载主机；虚拟机类型副本挂载对应的虚拟化平台
          </Alert>
          <Table searchable border :columns="columnsTitle"  @on-selection-change="selectionChange" :data="checkedData" stripe size="small">
            <template slot-scope="{ row }" slot="action">
              <div>
                <Button
                  type="success"
                  size="small"
                  style="margin-right: 5px"
                  @click="backupDetail(row)"
                >查看
                </Button>
                <Button
                  type="warning"
                  size="small"
                  style="margin-right: 5px"
                  @click="openEdit(row)"
                >编辑
                </Button>
                <Button v-if="row.status===1024"
                        type="info"
                        size="small"
                        style="margin-right: 5px"
                        @click="openMount(row)"
                >挂载</Button>
                <Poptip v-if="row.status===4096"
                        confirm
                        @on-ok="unMount(row)"
                        title="确定要卸载吗？">
                  <Button
                    type="error"
                    size="small"
                    style="margin-right: 5px"
                  >卸载
                  </Button>
                </Poptip>
                <Button v-if="row.status !==1024 && row.status !==4096"
                        type="info" disabled
                        size="small"
                        style="margin-right: 5px"
                        @click="openMount"
                >挂载</Button>
                <Button v-if="row.status !== 4096 && row.status !== 1024"
                  disabled
                  type="primary"
                  size="small"
                  style="margin-right: 5px"
                >快照</Button>
                <Button v-else
                        type="primary"
                        size="small"
                        style="margin-right: 5px"
                        @click="openSnapshot(row)"
                >快照</Button>
                <Button v-if="row.status===1024"
                  type='primary'
                  size="small"
                  style="margin-right: 5px"
                  @click="openMirror(row)"
                >镜像</Button>
                <Button v-else
                        disabled
                        type='primary'
                        size="small"
                        style="margin-right: 5px"
                >镜像</Button>
                <Poptip v-if="row.status===1024"
                  confirm
                  @on-ok="destroyBackup(row)"
                  title="确定要删除副本吗？">
                  <Button
                    type="error"
                    size="small"
                  >删除
                  </Button>
                </Poptip>
                <Poptip v-else
                    confirm
                    title="确定要删除副本吗？">
                  <Button
                    type="error" disabled
                    size="small"
                  >删除
                  </Button>
                </Poptip>
              </div>
            </template>
          </Table>
          <div class="table-page">
            <div class="table-page-position">
              <Page :total="total" :current="currentPage" show-total show-elevator show-sizer @on-change="changePage" @on-page-size-change="changeSize"></Page>
            </div>
          </div>
        </Row>
      </div>
      <Spin fix v-if="showSpin">
        <Icon type="ios-loading" size=18 class="demo-spin-icon-load" ></Icon>
        <div>Loading</div>
      </Spin>
    </Card>
    <snapshot
      :show="showCreateSnapshot"
      :backupName="currentBackup"
      :backupId="currentId"
      @closeCreate="closeCreateModal"
    />

    <Modal
      v-model="crtModal"
      draggable
      title="创建副本"
      @on-cancel="cancel">
      <Form :rules="ruleValidate" ref="formCrtPlat" :model="formCrtPlat">
        <FormItem label="名称" prop="name">
          <Input v-model="formCrtPlat.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formCrtPlat.desc" placeholder="请输入描述"/>
        </FormItem>
        <FormItem label="请选择数据类型" prop="type">
          <Select filterable v-model="formCrtPlat.type" >
            <Option v-for="(value, key) in typeMap" :value="key" :key="value" >{{ value }}</Option>
          </Select>
        </FormItem>
        <FormItem label="请选择所属池" prop="storePoolId">
          <Select filterable v-model="formCrtPlat.storePoolId" >
            <Option v-for="(value, key) in poolMap" :value="key" :key="value" >{{ value }}</Option>
          </Select>
        </FormItem>
        <FormItem label="大小(GB)" prop="size">
          <Input type="number" v-model="formCrtPlat.size" placeholder="请输入池大小"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick @click="crtPlat('formCrtPlat')">确定</Button>
      </div>
    </Modal>
    <Modal
      v-model="editModal"
      draggable
      title="编辑副本"
      @on-cancel="cancel">
      <Form  ref="formEditPlat" :model="formEditPlat">
        <FormItem label="名称" prop="desc">
          <Input v-model="formEditPlat.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formEditPlat.desc" placeholder="请输入描述"/>
        </FormItem>
        <strong style="color: #FF0033">
          只能扩大不能缩小
        </strong>
        <FormItem label="大小(GB)" prop="size">
          <Input type="number" v-model="formEditPlat.size" placeholder="请输入池大小"/>
        </FormItem>
        <strong style="color: #FF0033">
          修改副本类型后可能导致无法挂载，请谨慎修改！
        </strong>
        <FormItem label="请选择数据类型" prop="type">
          <Select filterable v-model="formEditPlat.type" >
            <Option v-for="(value, key) in typeMap" :value="key" :key="value" >{{ value }}</Option>
          </Select>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" v-preventReClick  @click="editBak('formEditPlat')">确定</Button>
      </div>
    </Modal>
    <Modal
      v-model="mountModal"
      title="挂载"
      width="1024px"
      :mask-closable = "false"
      @on-cancel="cancel">
      <div class="create-step" v-if="this.curFileType === 1004">
        <!-- 步骤条 -->
        <Steps :current="currentStep">
          <Step title="选择虚拟化平台"></Step>
          <Step title="选择位置路径"></Step>
          <Step title="选择主机"></Step>
          <Step title="选择计算资源"></Step>
          <Step title="完成"></Step>
        </Steps>
      </div>
      <div class="create-step" v-if=" this.curFileType === 1005 || this.curFileType === 1006">
        <!-- 步骤条 -->
        <Steps :current="currentStep">
          <Step title="选择虚拟化平台"></Step>
          <Step title="选择主机"></Step>
          <Step title="完成"></Step>
        </Steps>
      </div>
      <Row>
        <div>
          <Form :rules="ruleMountValidate" ref="formMount" :model="formMount">
            <div v-if="this.curFileType === 1003"><!--挂载主机-->
              <FormItem label="请选择挂载方式" >
                <RadioGroup v-model="mountTypeText">
                  <Radio v-for="item in mountTypeList" :label="item" >{{ item }}</Radio>
                </RadioGroup>
              </FormItem>
              <FormItem label="请选择主机类型" prop="platform">
                <el-select transfer="true" :popper-append-to-body="false" v-model="formMount.platform" placeholder="主机类型"  @change="getTargetList">
                  <el-option
                    v-for="(item, key) in platTypeList"
                    :key="key"
                    :label="item.key"
                    :value="item.value"></el-option>
                </el-select>
              </FormItem>
              <FormItem label="请选择挂载脚本"  v-if="this.mountTypeText==='应用默认挂载'">
                <el-select transfer="true" :popper-append-to-body="false" v-model="scriptName" placeholder="挂载脚本">
                  <el-option
                    v-for="(item, key) in scriptList"
                    :key="key"
                    :label="item.key"
                    :value="item.value"></el-option>
                </el-select>
              </FormItem>
            </div>

            <FormItem :label="'请选择'+platOrHost" prop="target" v-if="this.currentStep === 0">
              <el-select filterable transfer="true" :popper-append-to-body="false" v-model="formMount.target" :placeholder="platOrHost+'名称'">
                <el-option
                  v-for="(item, key) in mountTargetList"
                  :key="key"
                  :label="item.key"
                  :value="item.value">
                </el-option>
              </el-select>
            </FormItem>
            <!--------------------------CAS平台-------------------------->
            <div v-if="this.curFileType === 1005">
              <FormItem label="请选择主机" prop="casHostId" v-if="this.currentStep === 1">
                <el-select transfer="true" :popper-append-to-body="false" v-model="formMount.casHostInfo" placeholder="请选择主机" @change="selectHost">
                  <el-option
                    v-for="(item, key) in casHostList"
                    :key="key"
                    :label="item.key"
                    :value="item"></el-option>
                </el-select>
              </FormItem>
              <FormItem label="是否新建虚拟机"  v-if="this.currentStep === 2 ">
                <el-radio-group v-model="formMount.isCrtCas" @change="get_replica_info">
                  <el-radio label='false'>
                    <span>否(仅挂载NFS存储)</span>
                  </el-radio>
                  <el-radio label='true'>
                    <span>是</span>
                  </el-radio>
                </el-radio-group>
              </FormItem>
              <FormItem label="是否使用配置文件创建虚拟机" v-if="this.currentStep === 2 && formMount.isCrtCas==='true'">
                <el-radio-group v-model="formMount.isCrtCasByJson">
                  <el-radio-button label="false">否</el-radio-button>
                  <el-radio-button label="true">是</el-radio-button>
                </el-radio-group>
                  <strong style="color: #FF0033" v-if="this.currentStep === 2 && formMount.isCrtCas === 'true'">
                    使用默认配置创建虚拟机
                  </strong>
              </FormItem>
            </div>
            <!--------------------------FusionCompute平台-------------------------->
            <div v-if="this.curFileType === 1006">
              <FormItem label="请选择主机" prop="" v-if="this.currentStep === 1">
                <el-select transfer="true" :popper-append-to-body="false" v-model="formMount.fcHostInfo" placeholder="请选择主机" @change="selectHost">
                  <el-option
                    v-for="item in fcHostList"
                    :key="item.value"
                    :label="item.label"
                    :value="item"></el-option>
                </el-select>
              </FormItem>
              <FormItem label="是否新建虚拟机"  v-if="this.currentStep === 2 ">
                <el-radio-group v-model="formMount.isCrtFcVm">
                  <el-radio label='false'>
                    <span>否(仅挂载NFS存储)</span>
                  </el-radio>
                  <el-radio label='true'>
                    <span>是</span>
                  </el-radio>
                </el-radio-group>
              </FormItem>
              <FormItem label="是否立即刷新存储设备"  v-if="this.currentStep === 2 && this.formMount.isCrtFcVm ==='false'">
                <el-radio-group v-model="formMount.isRefresh">
                  <el-radio label='false'>
                    <span>否</span>
                  </el-radio>
                  <el-radio label='true'>
                    <span>是</span>
                  </el-radio>
                </el-radio-group>
                  <strong style="color: #FF0033" v-if="this.currentStep === 2 && this.formMount.isCrtFcVm ==='false'">
                    若不立即刷新，可以挂载多份存储后一次刷新
                  </strong>
              </FormItem>
              <FormItem label="操作系统" prop="os" v-if="this.curIdList.length===0 && this.currentStep === 2 && this.formMount.isCrtFcVm ==='true'">
                <el-select :popper-append-to-body="false" v-model="formMount.os" placeholder="请选择">
                  <el-option
                    v-for="item in osList"
                    :key="item"
                    :value="item">
                  </el-option>
                </el-select>
              </FormItem>
            </div>
            <!--------------------------VMware平台-------------------------->
              <FormItem label="请选择位置路径" v-if="this.currentStep === 1 && this.curFileType === 1004">
                <RadioGroup v-model="formMount.locationPath" vertical>
                  <Radio v-for="item in locationPathList" :label="item" >{{ item }}</Radio>
                </RadioGroup>
              </FormItem>
              <FormItem label="请选择主机" v-if="this.currentStep === 2 && this.curFileType === 1004">
                <RadioGroup v-model="formMount.vmwareHost" vertical>
                  <Radio v-for="item in vmwareHostList" :label="item" >{{ item }}</Radio>
                </RadioGroup>
                <strong style="color: #FF0033" v-if="this.currentStep === 2  && this.curFileType === 1004 && vmwareHostFlag === false">
                  该位置下没有主机,请选择其他位置
                </strong>
              </FormItem>
              <Tree :data="resourcePathList"
                v-if="this.currentStep === 3 && this.curFileType === 1004"
                :show-checkbox="true"
                :check-strictly="true"
                :check-directly="true"
                @on-check-change='selectLocationPath'>
              </Tree>
              <div v-if="this.currentStep === 4 && this.curFileType === 1004">
                <FormItem label="是否创建虚拟机">
                  <el-radio-group v-model="formMount.isCustom">
                    <el-radio label="false">否</el-radio>
                    <el-radio label="true">是</el-radio>
                  </el-radio-group>
                    <strong style="color: #FF0033" v-if="formMount.isCustom === 'true' && this.curIdList.length===0">
                      注册虚拟机，副本中必须包含有效的虚拟机磁盘文件，否则无法注册成功
                    </strong>
                </FormItem>
                <div v-if="formMount.isCustom === 'true'">
                  <FormItem label="是否打开虚拟机电源" prop="poweron" v-if="this.curIdList.length===0">
                    <el-radio-group v-model="formMount.poweron">
                      <el-radio label="false">否</el-radio>
                      <el-radio label="true">是</el-radio>
                    </el-radio-group>
                  </FormItem>
                  <FormItem label="操作系统" prop="os" v-if="this.curIdList.length===0">
                    <el-select  v-model="formMount.os" placeholder="请选择">
                      <el-option
                        v-for="item in osList"
                        :key="item"
                        :value="item">
                      </el-option>
                    </el-select>
                  </FormItem>
                  <FormItem label="虚拟机名称" prop="vmname" v-if="this.curIdList.length===0">
                    <el-input v-model="formMount.vmname" placeholder="请输入虚拟机名称" clearable></el-input>
                    <strong style="color: #FF0033" v-if="this.curIdList.length===0">
                      默认使用副本名称
                    </strong>
                  </FormItem>
                  <FormItem label="用户名" prop="username" v-if="this.curIdList.length===0">
                    <el-input v-model="formMount.vmUsername" placeholder="请输入虚拟机用户名"  clearable></el-input>
                  </FormItem>
                  <FormItem label="密码" prop="password" v-if="this.curIdList.length===0">
                    <el-input v-model="formMount.vmPassword" placeholder="请输入虚拟机密码"  clearable></el-input>
                  </FormItem>
                  <FormItem label="是否配置网络" prop="isSetIp" v-if="this.curIdList.length===0">
                    <el-radio-group v-model="formMount.isSetIp">
                      <el-radio label="false">否</el-radio>
                      <el-radio label="true">是</el-radio>
                    </el-radio-group>
                  </FormItem>
                  <div v-if="formMount.isSetIp === 'true'">
                    <FormItem label="IP地址" prop="addr">
                      <el-input v-model="formMount.addr" placeholder="请输入虚拟机IP地址"  clearable></el-input>
                    </FormItem>
                    <FormItem label="网关" prop="gateway">
                      <el-input v-model="formMount.gateway" placeholder="请输入虚拟机网关" clearable></el-input>
                    </FormItem>
                    <FormItem label="子网掩码" prop="netmask">
                      <el-input v-model="formMount.netmask" placeholder="请输入子网掩码" clearable></el-input>
                    </FormItem>
                  </div>
                </div>
              </div>
          </Form>
<!--          <Upload ref="upload"-->
<!--            type="drag"-->
<!--            style="width: 200px"-->
<!--            :before-upload="beforeUpload"-->
<!--            :on-success="successUpload"-->
<!--            :on-error="failUpload"-->
<!--            :on-preview="clickFile"-->
<!--            :action=upUrl-->
<!--            :headers="upHeaders"-->
<!--            :max-size="10000"-->
<!--            v-if="this.curFileType===1003 && this.mountTypeText==='应用默认挂载'"-->
<!--          >-->
<!--            <div style="padding: 20px 0">-->
<!--              <Icon type="ios-cloud-upload" size="22" style="color: #3399ff"></Icon>-->
<!--              <p>上传挂载脚本(小于10M)</p>-->
<!--            </div>-->
<!--          </Upload>-->
        </div>
      </Row>
      <div slot="footer"><!-- 底部跳转按钮-->
        <!--步骤①-->
        <div v-if="this.currentStep === 0">
          <Button type="primary" v-preventReClick :disabled="isdisabledFn" @click="toNextStep1('formMount')"
                  v-if="this.currentStep === 0 && (this.curFileType === 1004 ||this.curFileType === 1005 ||this.curFileType === 1006 )">下一步
          </Button>
        </div>

        <!--步骤二-->
        <div v-if="this.currentStep === 1">
          <Button type="primary"  @click="changeCurrentStep(0)" >上一步</Button>
          <Button type="primary"  @click="toNextStep2_" :disabled="isdisabledFn" v-if="this.curFileType===1004 ">下一步</Button>
          <Button type="primary"  @click="toNextStep2Cas('formMount')" :disabled="isdisabledFn" v-if="this.curFileType===1005">下一步</Button>
          <Button type="primary"  @click="toNextStep2Fc('formMount')" :disabled="isdisabledFn" v-if="this.curFileType===1006">下一步</Button>
        </div>

        <!--步骤三-->
        <div v-if="this.currentStep === 2">
          <Button type="primary"  @click="changeCurrentStep(1)" >上一步</Button>
          <Button type="primary"  @click="toNextStep3('formMount')"  :disabled="isdisabledFn" v-if="this.curFileType===1004 ">下一步</Button>
          <Button type="primary" v-preventReClick @click="doMountCAS('formMount')" v-if="this.curFileType===1005">确定</Button>
          <Button type="primary" v-preventReClick @click="doMountFc('formMount')" v-if="this.curFileType===1006">确定</Button>
        </div>
        <!--步骤四-->
        <div v-if="this.currentStep === 3">
          <Button type="primary"  @click="changeCurrentStep(2)" v-if="this.curFileType===1004">上一步</Button>
          <Button type="primary"  @click="toNextStep4('formMount')"  v-if="this.curFileType===1004 ">下一步</Button>
          <Button type="primary"  @click="changeCurrentStep(0)" v-if="this.formMount.platform==='51'">上一步</Button>
        </div>

        <!--步骤五-->
        <div v-if="this.currentStep === 4">
          <Button type="primary"  @click="changeCurrentStep(3)" v-if="this.curFileType===1004">上一步</Button>
          <Button type="primary"  @click="changeCurrentStep(0)" v-if="this.formMount.platform==='51'">上一步</Button>
          <Button type="primary" v-preventReClick @click="doMountVM('formMount')" v-if="this.curFileType===1004">确定</Button>
        </div>

        <Button type="primary" v-preventReClick @click="doMountCus('formMount')" v-if="this.curFileType === 1003">确定</Button>

      </div>
    </Modal>
    <Modal
      draggable
      v-model="mirrorModal"
      title="生成镜像"
      @on-cancel="cancel">
      <Form :rules="ruleMirrorValidate" ref="formMirror" :model="formMirror">
        <FormItem label="名称" prop="name">
          <Input v-model="formMirror.name" placeholder="请输入名称"/>
        </FormItem>
        <FormItem label="描述" prop="desc">
          <Input v-model="formMirror.desc" placeholder="请输入描述"/>
        </FormItem>
      </Form>
      <div slot="footer">
        <Button type="primary" :disabled="isdisabledFn" v-preventReClick @click="crtMirror('formMirror')">确定</Button>
      </div>
    </Modal>
    <Modal title="副本信息" width="1400px" footer-hide v-model="showDetail">
      <backup-info :b-id=curBuildId :typeMap="typeMap" ref="backupInfo"></backup-info>
    </Modal>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import {
  getTypes,
  getAllPools,
  crtBackup,
  getBackupList,
  getImageStatus,
  choiceImageStateColor,
  toGB,
  crtMirror,
  deleteBackup,
  mountBackup,
  unmountBackup,
  getResource,
  getMountType,
  editBackup,
  mountBackupBatch,
  getAllScriptList,
  unMountBackupBatch,
  deleteBackupBatch,
  getCasHosts,
  choiceTypeColor,
  getVMwareHostsByPath,
  getCasHostInfo, backupDetails, getFcHosts
} from '@/api/backup.js'
import snapshot from './snapshot'
import { httpRequest } from '@/libs/tools'
import backupInfo from '_c/deploy/backup-info'
import { getHostListAll, getHostTypes, getPlatformTypes, getPlatListAll} from '@/api/env'
import {getToken} from '@/libs/util'
import config from '@/config'
import JFSocketIo, {sendWS} from '@/libs/jf-socket-io'
const baseUrl = process.env.NODE_ENV === 'development' ? config.baseUrl.dev : config.baseUrl.pro

export default {
  name: 'backupList',
  components: {
    snapshot,
    backupInfo
  },
  // 子组件自己的数据定义，可读可写
  data () {
    return {
      // socketIo: new JFSocketIo('replica'),
      backId: '',
      showDetail: false,
      isdisabledFn: false,
      platOrHost: '',
      curBuildId: '',
      curIdList: [],
      curSize: '',
      curFileType: '',
      replicaName: '',
      isExecScript: '',
      scriptName: '',
      editModal: false,
      // 搜索
      searchObj: {
        backupName: '',
        status: '',
        dataType: '',
        uuid: ''
      },
      formCrtPlat: {
        type: '',
        name: '',
        desc: '',
        size: '',
        storePoolId: ''
      },
      formEditPlat: {
        name: '',
        desc: '',
        size: '',
        type: ''
      },
      formMount: {
        type: '',
        platform: '',
        target: '',
        locationPath: '',
        resourcePath: '',
        vmwareHost: '',
        isCustom: 'false',
        vmname: '',
        hostname: '',
        isSetIp: '',
        addr: '',
        netmask: '',
        gateway: '',
        os: 'Linux',
        vmUsername: '',
        vmPassword: '',
        poweron: '',
        casHostId: '',
        casHostName: '',
        casHostInfo: {},
        casClusterId: '',
        casHostPoolId: '',
        casVmId: '',
        casVmName: '',
        casStoreFile: '',
        isCrtCas: '',
        isCrtCasByJson: '',
        fcHostInfo: {},
        isCrtFcVm: '',
        isRefresh: '',
        replicaName: '',
        poolName: ''
        // mountPoint:''
      },
      formMirror: {
        backupId: '',
        name: '',
        desc: ''
      },
      typeMap: {},
      typeMapQuery: { 0: '全部' },
      platTypeMap: {},
      platTypeList: [],
      scriptList:[],
      mountTypeList:[],
      poolMap: {},
      crtModal: false,
      mountModal: false,
      mirrorModal: false,
      currentStep: 0,
      dataTypeList: ['主机', '虚拟化平台'],
      osList: ['Linux', 'Windows'],
      mountTargetList: [],
      mountType: '',
      mountTypeText: '应用默认挂载',
      mountTypeMap: {},
      locationPathList: [],
      resourcePathList: [],
      vmwareClusterList: [],
      vmwareHostList: [],
      vmwareHostFlag: false,
      casHostList: [],
      casVmList: [],
      casStorageList: [],
      fcHostList: [],
      selectBox: [],
      showSpin: false,
      upHeaders: {
        Authorization: 'Tangula ' + getToken()
      },
      upUrl: config.baseUrl.pro + '/file/upload',
      ruleValidate: {
        // reviewConclusion: [{ required: true, message: '请输入', trigger: 'change' }],
        // lastDate: [{ required: true, message: '请输入', trigger: 'change', type: 'date' }],
        name: [{ required: true, message: '请输入', trigger: 'blur' }],
        type: [{ required: true, message: '请选择', trigger: 'blur' }],
        size: [{ required: true, message: '请输入', trigger: 'blur' }],
        storePoolId: [{ required: true, message: '请选择', trigger: 'blur' }],
        os: [{ required: true, message: '请选择', trigger: 'blur' }],
      },
      ruleMirrorValidate: {
        name: [{ required: true, message: '请输入', trigger: 'blur' }]
      },
      ruleMountValidate: {
        type: [{ required: true, message: '请选择', trigger: 'blur' }],
        platform: [{ required: true, message: '请选择', trigger: 'blur' }],
        os: [{ required: true, message: '请选择', trigger: 'blur' }],
        addr: [{ required: true, message: '输入', trigger: 'blur' }],
        netmask: [{ required: true, message: '输入', trigger: 'blur' }],
        gateway: [{ required: true, message: '输入', trigger: 'blur' }],
        target: [{ required: true, message: '请选择', trigger: 'blur' }],
        // casHostId: [{ required: true, message: '请选择', trigger: 'blur' }],
        casVmId: [{ required: true, message: '请选择', trigger: 'blur' }],
        casStoreFile: [{ required: true, message: '请选择', trigger: 'blur' }],
        casVmName: [{ required: true, message: '输入', trigger: 'blur' }]
      },
      currentId:0,
      currentBackup: '',
      searchChanged: false,
      backupStatus: [
        { key: '全部', value: 0 },
        { key: '未挂载', value: 1024 },
        { key: '挂载中', value: 2048 },
        { key: '已挂载', value: 4096 },
        { key: '卸载中', value: 8192 }
      ],
      // 需求详情模态框
      total: 0,
      currentPage: 1,
      pageSize: 10,
      orders: [],
      showCreateSnapshot: false,
      columnsTitle: [
        {
          type: 'selection',
          selectable: 'checkbox',
          width: 55,
          align: 'left'
        },
        {
          title: 'ID',
          key: 'id',
          width: 65,
          align: 'left'
        },
        {
          title: '名称',
          key: 'name',
          align: 'center',
          resizable: true,
          width: 200,
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
          title: '大小(GB)',
          width: 90,
          key: '',
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
          width: 125,
          key: 'type',
          align: 'center',
          render: (h, params) => {
            let color = choiceTypeColor(params.row.type)
            return h('span', {
              style: {
                color: color
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
              }, params.row.status)
            }
          }
        },
        // {
        //   title: '导出路径',
        //   key: 'name',
        //   align: 'center',
        //   width: 150,
        //   resizable: true,
        //   render: (h, params) => {
        //     let arr = []
        //     arr.push(
        //       h('span', {
        //           style: { cursor: 'pointer' }
        //         },
        //         [
        //           h('span',
        //             {
        //               style: {
        //                 display: 'inline-block',
        //                 width: '100%',
        //                 overflow: 'hidden',
        //                 textOverflow: 'ellipsis',
        //                 whiteSpace: 'nowrap'
        //               },
        //               domProps: {
        //                 title: params.row.export
        //               }
        //             },
        //             params.row.export)
        //         ]
        //       )
        //
        //     )
        //     return h('span', arr)
        //   }
        // },
        {
          title: '创建时间',
          width: 100,
          key: 'created_time',
          align: 'center'
        },

        {
          title: '操作',
          key: '',
          width: 305,
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
    }),
    checkedData (state) {
      return state.orders.map(d => {
        // if(d.status!==1024){
        //   d['_disabled'] = true
        // }
        return d
      })
    }
  },
  mounted: function () {
    this.getTypes()
    this.getPools()
    this.getData()
    this.getWebsocket()
    // this.socketIo.sendWS('hello websockets')
    // this.socketIo.websock.onmessage = (event) => {
    //   if(String(event.data)===this.userName){
    //     this.getData()
    //     console.log('websocket消息推送', event.data)
    //   }
    // }
  },
  methods: {
    getWebsocket(){
      let newSocket=new JFSocketIo('replica')
      newSocket.sendWS('hello websockets')

      // 接收消息
      const getsocketData = e => { // 创建接收消息函数
        const data = e && e.detail.data
        if (String(data) === this.userName) {
          console.log('websocket消息推送', data)
          this.getData(this.currentPage)
        }
      }
      // 注册监听事件
      window.addEventListener('onmessageWS', getsocketData)

      // window.onbeforeunload = function(){
      //   console.log('页面刷新')
      //   newSocket.closeConnection()
      // }
    },
    selectionChange (selection) {
      this.selectBox = selection
    },
    openCreate () {
      this.crtModal = true
    },
    changeCurrentStep (step) {
      this.currentStep = step
    },
    getData (page = 1, pageSize = this.pageSize) {
      this.showSpin = true
      this.selectBox = []
      const params = Object.assign({}, {
        status: this.searchObj.status,
        filter: this.searchObj.backupName,
        type: this.searchObj.dataType,
        uuid: this.searchObj.uuid
      })

      httpRequest(this, getBackupList, [page, pageSize, params]).then(res => {
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
    // 副本类型
    getTypes () {
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
    // 挂载类型
    getMountTypes () {
      httpRequest(this, getMountType, []).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        this.mountTypeList=[]
        this.mountTypeMap={}
        for(let key in res.response){
          let value=res.response[key]
          this.mountTypeList.push(value)
          this.mountTypeMap[value]=key
        }
      })
    },
    // 平台类型列表
    getTypeList(type){
      this.platTypeList = []
      if(type === '主机'){
        httpRequest(this, getHostTypes, []).then(res => {
          if (res.code !== 200) {
            this.$Notice.error({
              title: '错误',
              desc: res.message
            })
            return
          }
          for (let key in res.response){
            let tmpMap ={key:res.response[key],value:key}
            this.platTypeList.push(tmpMap)
          }
        })
      }else if(type === '虚拟化平台'){
        httpRequest(this, getPlatformTypes, []).then(res => {
          if (res.code !== 200) {
            this.$Notice.error({
              title: '错误',
              desc: res.message
            })
            return
          }
          for (let key in res.response){
            let tmpMap ={key:res.response[key],value:key}
            this.platTypeList.push(tmpMap)
          }
        })
      }
    },
    // 主机/平台列表
    getTargetList () {
      this.mountTargetList = []
      if (this.curFileType === 1004 || this.curFileType === 1005 || this.curFileType === 1006) { //虚拟化平台
        const params = Object.assign({}, {
          type: this.formMount.platform
        })
        httpRequest(this, getPlatListAll, [params]).then(res => {
          this.showSpin = false
          if (res.code !== 200) {
            this.$Notice.error({
              title: '错误',
              desc: res.message
            })
            return
          }
          if (res.response.data) {
            for(let i =0;i<res.response.data.length;i++){
              this.mountTargetList.push({key:res.response.data[i].name,value:String(res.response.data[i].id)}) //String用于过参数校验
            }
          }
        })
      }else {
        const params = Object.assign({},{
          type: this.formMount.platform,
        })
        httpRequest(this, getHostListAll, [params]).then(res => {
          this.showSpin = false
          if (res.code !== 200) {
            this.$Notice.error({
              title: '错误',
              desc: res.message
            })
            return
          }
          if(res.response.data){
            for(let i =0;i<res.response.data.length;i++){
              this.mountTargetList.push({key:res.response.data[i].name,value:String(res.response.data[i].id)})
            }
          }
        })
      }
    },
    getPools () {
      const params = Object.assign({}, {
        // filter: this.searchObj.poolName,
      })
      httpRequest(this, getAllPools, []).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '获取存储池失败',
            desc: res.message
          })
          return
        }
        if(res.response.data){
          for(var i=0;i<res.response.data.length;i++){
            this.poolMap[res.response.data[i].uuid] = res.response.data[i].name
          }
        }
      })
    },

    crtPlat (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          const param = {
            name: this.formCrtPlat.name,
            desc: this.formCrtPlat.desc,
            type: Number(this.formCrtPlat.type),
            size: Number(this.formCrtPlat.size),
            storePoolId: this.formCrtPlat.storePoolId
          }
          crtBackup(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功'
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
    editBak (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          if ( Number(this.formEditPlat.size) < Number(this.curSize)) {
            this.$Message.error('修改后的副本大小不能小于当前')
            return
          }
          let type
          if (!this.formEditPlat.type) {
            type = 0
          } else {
            type = Number(this.formEditPlat.type)
          }
          const param = {
            id: this.curBuildId,
            size: Number(this.formEditPlat.size),
            name: this.formEditPlat.name,
            desc: this.formEditPlat.desc,
            type: type
          }
          editBackup(param).then(({ data: { code, message } }) => {
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
            this.getData(this.currentPage)
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
    openEdit (row) {
      this.curBuildId = row.id
      this.curSize = toGB(row.size)
      this.formEditPlat.desc = row.desc
      this.formEditPlat.size = toGB(row.size)
      this.editModal = true
    },
    openMount (row) {
      this.currentStep = 0
      this.curBuildId = row.id
      this.curFileType = row.type
      this.platTypeList = []
      this.mountTargetList = []
      this.formMount = {}
      this.$set(this.formMount, 'isCustom', 'false')
      this.$set(this.formMount, 'poweron', 'false')
      this.$set(this.formMount, 'isSetIp', 'false')
      this.$set(this.formMount, 'isCrtCas', 'false')
      this.$set(this.formMount, 'isCrtCasByJson', 'false')
      this.$set(this.formMount, 'replicaName', row.name)
      this.$set(this.formMount, 'poolName', row.uuid)
      this.$set(this.formMount, 'isCrtFcVm', 'false')
      this.$set(this.formMount, 'isRefresh', 'true')
      this.scriptName = ''

      if (row.type === 1004) { // VMware类型
        this.formMount.platform = 11
        this.platOrHost = '虚拟化平台'
        this.getTargetList()
      } else if (row.type === 1005) { // CAS类型
        this.formMount.platform = 12
        this.platOrHost = '虚拟化平台'
        this.getTargetList()
      } else if (row.type === 1006) { // FusionCompute类型
        this.formMount.platform = 13
        this.platOrHost = '虚拟化平台'
        this.getTargetList()
      } else {
        this.getTypeList('主机')
        this.platOrHost = '主机'
        this.getMountTypes()
        this.getScriptList()
      }
      this.mountModal = true
    },
    printIsCustom () {
      console.log(this.formMount.isCustom)
    },
    openMountBatch () {
      if (this.selectBox.length === 0) {
        this.$Message.error('请选择副本')
        return
      }
      let tmpList = []
      for (let i in this.selectBox) {
        if (this.selectBox[i].status !== 1024) {
          this.$Message.error('只有未挂载状态的副本可以挂载')
          return
        }
        tmpList.push(this.selectBox[i].type)
      }
      if (tmpList.length > 0 && (new Set(tmpList)).size !== 1) {
        this.$Message.error('请选择相同数据类型的副本进行批量挂载')
        return
      }
      this.curFileType = tmpList[0]
      this.mountModal = true
      this.currentStep = 0
      for (let i in this.selectBox) {
        this.curIdList.push(this.selectBox[i].id)
      }
      this.formMount = {}
      this.$set(this.formMount, 'isCustom', 'false')
      this.$set(this.formMount, 'poweron', 'false')
      this.$set(this.formMount, 'isSetIp', 'false')
      this.$set(this.formMount, 'isCrtCas', 'false')
      this.$set(this.formMount, 'isCrtCasByJson', 'false')

      this.platTypeList = []
      this.mountTargetList = []
      this.scriptName = ''
      if (this.curFileType === 1004) { // VMware类型
        this.formMount.platform = 11
        this.getTargetList()
      } else if (this.curFileType === 1005) { // CAS类型
        this.formMount.platform = 12
        this.getTargetList()
      } else if (row.type === 1006) { // FusionCompute类型
        this.formMount.platform = 13
        this.platOrHost = '虚拟化平台'
        this.getTargetList()
      } else {
        this.getTypeList('主机')
        this.getMountTypes()
        this.getScriptList()
      }
    },
    batchUnMount () {
      if (this.selectBox.length === 0) {
        this.$Message.error('请选择副本')
        return
      }
      let replicaIdList = []
      for (let i in this.selectBox) {
        if (this.selectBox[i].status !== 4096) {
          this.$Message.error('只有已挂载的副本可以卸载')
          return
        }
        replicaIdList.push(this.selectBox[i].id)
      }
      this.$Modal.confirm({
        title: '确认卸载？',
        onOk: () => {
          this.$Notice.info({
            title: '请求中……',
            desc: '正在批量卸载，请等待'
          })
          const param = {
            unMountInfo: replicaIdList
          }
          unMountBackupBatch(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功,请稍后查看卸载结果'
            })
            this.mountModal = false
            this.getData(this.currentPage)
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
    batchDelete(){
      if(this.selectBox.length===0){
        this.$Message.error('请选择副本')
        return
      }
      let replicaIdList=[]
      for(let i in this.selectBox){
        if(this.selectBox[i].status!==1024){
          this.$Message.error('只有空闲的副本可以删除')
          return
        }
        replicaIdList.push(this.selectBox[i].id)
      }
      this.$Modal.confirm({
        title: '确认删除？',
        onOk: () => {
          this.$Notice.info({
            title: '请求中……',
            desc: '正在批量删除，请等待'
          })
          const param = {
            replicaIdList: replicaIdList
          }
          deleteBackupBatch(param).then(({data: {code, message}}) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功,请稍后查看卸载结果'
            })
            this.mountModal = false
            this.getData(this.currentPage)
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
    toNextStep1 (name) {
      this.locationPathList = []
      this.$refs[name].validate((valid) => {
        if (valid) {
          if (this.curFileType === 1004) { // vmware
            this.isdisabledFn = true
            const params = Object.assign({}, {
              id: Number(this.formMount.target),
              showType: 2,
              isGetAll: 'false'
            })
            this.$Notice.info({
              title: '请求中……',
              desc: '正在查询数据源,请稍候'
            })
            httpRequest(this, getResource, [params]).then(res => {
              this.formMount.locationPath = ''
              if (res.code !== 200) {
                this.$Notice.error({
                  title: '获取数据源失败 ',
                  desc: res.message
                })
                this.isdisabledFn = false
                return
              }
              this.$Notice.success({
                title: '成功',
                desc: '请求成功'
              })
              this.isdisabledFn = false
              this.convertList(res.response.data)
              this.changeCurrentStep(1)
            })
          } else if (this.curFileType === 1005) { // CAS
            this.changeCurrentStep(1)
            this.isdisabledFn = true
            this.getCasHostList()
            this.isdisabledFn = false
          } else if (this.curFileType === 1006) { // fusioncompute
            this.changeCurrentStep(1)
            this.isdisabledFn = true
            this.getFcHostList()
            this.isdisabledFn = false
          } else {
            this.changeCurrentStep(3)
          }
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    toNextStep2_ () {
      this.vmwareHostList = []
      if (this.curFileType === 1004 && !this.formMount.locationPath) {
        this.$Message.error('位置路径不能为空!')
      } else {
        this.isdisabledFn = true
        this.formMount.vmwareHost = ''
        const params = Object.assign({}, {
          id: Number(this.formMount.target),
          path: this.formMount.locationPath
        })
        this.$Notice.info({
          title: '请求中……',
          desc: '正在查询数据源,请稍候'
        })
        httpRequest(this, getVMwareHostsByPath, [params]).then(res => {
          if (res.code !== 200) {
            this.$Notice.error({
              title: '获取数据源失败',
              desc: res.message
            })
            this.isdisabledFn = false
            return
          }
          this.$Notice.success({
            title: '成功',
            desc: '请求成功'
          })
          this.isdisabledFn = false
          console.log('查询主机列表返回：', res.response.data)
          if (res.response.data) {
            this.vmwareHostList = res.response.data
            this.vmwareHostFlag = true
          } else {
            this.vmwareHostFlag = false
          }
          this.changeCurrentStep(2)
        })
      }
    },
    toNextStep2Fc (name) {
      if (!this.formMount.fcHostInfo) {
        this.$Message.error('请选择主机!')
        return
      }
      /**
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn = true
          this.formMount.casHostPoolId = ''
          this.formMount.casClusterId = ''
          this.$Notice.info({
            title: '请求中……',
            desc: '正在查询,请稍候'
          })
          const params = Object.assign({}, {
            id: Number(this.formMount.target),
            // hostId: Number(this.formMount.casHostId)
            hostId: Number(this.formMount.casHostInfo.value)
          })
          httpRequest(this, getCasHostInfo, [params]).then(res => {
            if (res.code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: res.message
              })
              this.isdisabledFn = false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功'
            })
            setTimeout(() => {
              if (res.response) {
                this.formMount.casHostPoolId = res.response.hostPoolId
                this.formMount.casClusterId = res.response.clusterId
                console.log('查询主机信息:',res.response)
              } else {
                console.log('查询主机信息为空')
              }
              this.isdisabledFn = false
              this.changeCurrentStep(2)
            }, 300)
          })
        } else {
          this.$Message.error('请输入必填项!')
        }

      })
       **/
      this.changeCurrentStep(2)
    },
    toNextStep2Cas (name) {
      if (!this.formMount.casHostInfo) {
        this.$Message.error('请选择主机!')
        return
      }
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn = true
          this.formMount.casHostPoolId = ''
          this.formMount.casClusterId = ''
          this.$Notice.info({
            title: '请求中……',
            desc: '正在查询,请稍候'
          })
          const params = Object.assign({}, {
            id: Number(this.formMount.target),
            // hostId: Number(this.formMount.casHostId)
            hostId: Number(this.formMount.casHostInfo.value)
          })
          httpRequest(this, getCasHostInfo, [params]).then(res => {
            if (res.code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: res.message
              })
              this.isdisabledFn = false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功'
            })
            setTimeout(() => {
              if (res.response) {
                this.formMount.casHostPoolId = res.response.hostPoolId
                this.formMount.casClusterId = res.response.clusterId
                console.log('查询主机信息:',res.response)
              } else {
                console.log('查询主机信息为空')
              }
              this.isdisabledFn = false
              this.changeCurrentStep(2)
            }, 300)
          })
          // this.changeCurrentStep(2)
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    toNextStep3 () {
      this.resourcePathList = []
      if (this.curFileType === 1004 && !this.formMount.vmwareHost) {
        this.$Message.error('主机不能为空!')
      } else {
        this.isdisabledFn = true
        this.formMount.resourcePath = ''
        const params = Object.assign({}, {
          id: Number(this.formMount.target),
          fullPath: this.formMount.locationPath,
          hostName: this.formMount.vmwareHost,
          showType: 1,
          isGetAll: 'true'
        })
        this.$Notice.info({
          title: '请求中……',
          desc: '正在查询数据源,请稍候'
        })
        httpRequest(this, getResource, [params]).then(res => {
          if (res.code !== 200) {
            this.$Notice.error({
              title: '获取数据源失败',
              desc: res.message
            })
            this.isdisabledFn = false
            return
          }
          this.$Notice.success({
            title: '成功',
            desc: '请求成功'
          })
          this.isdisabledFn = false
          if (res.response.data) {
            this.resourcePathList = this.convertTree(res.response.data)
          }
          this.changeCurrentStep(3)
        })
      }
    },
    toNextStep4 () {
      if (this.curFileType === 1004 && !this.formMount.resourcePath) {
        this.$Message.error('计算资源不能为空')
      } else {
        this.changeCurrentStep(4)
      }
    },
    selectLocationPath (arr,obj) {
      // console.log('选中树:', obj)
      this.formMount.resourcePath = obj.path
      arr.forEach(item => {
        item.checked = false
      })
      // 只选中最后一次选中的
      obj.checked = true
      console.log('计算资源:', this.formMount.resourcePath)
    },

    convertList (tree) {
      const result = []
      tree.forEach((item) => {
        // 读取 map 的键值映射
        let expand = false
        let title = item['path']
        let children = item['subObject']
        if (item['type'] === 'Datacenter') {
          this.locationPathList.push(item['path'])
        }
        // 如果有子节点，递归
        if (children) {
          this.convertList(children)
        }
        result.push({ expand, title, children, ...item })
      })
    },

    convertTree (tree) {
      const result = []
      tree.forEach((item) => {
        // 读取 map 的键值映射
        let expand = item['expanded']
        let title = item['name']
        let children = item['subObject']
        // 如果有子节点，递归
        if (children) {
          children = this.convertTree(children)
        }
        result.push({ expand, title, children, ...item })
      })
      return result
    },
    doMountVM (name) { // 挂载VMware类型副本
      this.$refs[name].validate((valid) => {
        if (valid) {
          let mountInfo = {
            replicaId: 0,
            targetType: Number(this.formMount.platform),
            targetId: Number(this.formMount.target),
            mountType: 1
          }
          if (!this.formMount.isCustom) {
            this.$Message.error('请选择是否注册虚拟机')
            return
          }
          Object.assign(
            mountInfo, {
              isExecScript: false
            })
          let appConfig = {
            locationPath: this.formMount.locationPath,
            computeResource: this.formMount.resourcePath,
            vmwareHost: this.formMount.vmwareHost
          }
          if (this.formMount.isCustom === 'true') { // 手动注册虚拟机
            let isSetIp
            isSetIp = this.formMount.isSetIp === 'true'
            Object.assign(
              appConfig, {
                isRegisterVM: true,
                powerOn: this.formMount.poweron,
                vmName: this.formMount.vmname,
                username: this.formMount.vmUsername,
                password: this.formMount.vmPassword,
                isSetIp: isSetIp,
                addr: this.formMount.addr,
                netMask: this.formMount.netmask,
                gateWay: this.formMount.gateway,
                os: this.formMount.os
              })
          } else { // 否
            Object.assign(
              appConfig, {
                isRegisterVM: false
              })
          }
          Object.assign(
            mountInfo, {
              appConfig: appConfig
            })
          // 单个副本挂载
          if (this.curIdList.length === 0) {
            mountInfo.replicaId = this.curBuildId
            this.$Modal.confirm({
              title: '确认挂载？',
              onOk: () => {
                this.$Notice.info({
                  title: '请求中……',
                  desc: '正在挂载，请等待'
                })
                const param = {
                  mountInfo: mountInfo
                }
                mountBackup(param).then(({ data: { code, message } }) => {
                  if (code !== 200) {
                    this.$Notice.error({
                      title: '错误',
                      desc: message
                    })
                    return
                  }
                  this.$Notice.success({
                    title: '成功',
                    desc: '请求成功,请稍后查看挂载结果'
                  })
                  this.mountModal = false
                  this.getData(this.currentPage)
                }).catch(err => {
                  if (err && err.response) {
                    this.$Notice.error({
                      title: '错误',
                      desc: err.response.data.msg
                    })
                  }
                })
              },
              onCancel: () => {
              }
            })
          } else { // 批量挂载
            let mountInfoList = []
            for (let i in this.curIdList) {
              var tmpMap = Object.assign({}, mountInfo)
              tmpMap.replicaId = this.curIdList[i]
              mountInfoList.push(tmpMap)
            }
            this.$Modal.confirm({
              title: '确认挂载？',
              onOk: () => {
                this.$Notice.info({
                  title: '请求中……',
                  desc: '正在批量挂载，请等待'
                })
                const param = {
                  mountInfos: mountInfoList
                }
                mountBackupBatch(param).then(({ data: { code, message } }) => {
                  if (code !== 200) {
                    this.$Notice.error({
                      title: '错误',
                      desc: message
                    })
                    return
                  }
                  this.$Notice.success({
                    title: '成功',
                    desc: '请求成功,请稍后查看挂载结果'
                  })
                  this.curIdList = []
                  this.mountModal = false
                  this.getData(this.currentPage)
                }).catch(err => {
                  if (err && err.response) {
                    this.$Notice.error({
                      title: '错误',
                      desc: err.response.data.msg
                    })
                  }
                })
              },
              onCancel: () => {
              }
            })
          }
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },

    doMountCAS (name) { // 挂载CAS平台
      let mountInfo = {
        replicaId: 0,
        targetType: Number(this.formMount.platform),
        targetId: Number(this.formMount.target),
        mountType: 1
      }
      if (!this.formMount.isCrtCas) {
        this.$Message.error('请选择是否创建虚拟机')
        return
      }
      Object.assign(
        mountInfo, {
          isExecScript: false
        })
      let appConfig = {
        // hostId: this.formMount.casHostId,
        hostId: Number(this.formMount.casHostInfo.value),
        hostName: this.formMount.casHostInfo.key,
        storeName: ''
      }
      if (this.formMount.isCrtCas === 'true') { // 注册虚拟机
        let isCrtCasByJson
        isCrtCasByJson = this.formMount.isCrtCasByJson === 'true'
        Object.assign(
          appConfig, {
            isRegisterVM: true,
            isCrtCasByJson: isCrtCasByJson,
            clusterId: this.formMount.casClusterId,
            hostPoolId: this.formMount.casHostPoolId
          })
      } else { // 否
        Object.assign(
          appConfig, {
            isRegisterVM: false
          })
      }
      Object.assign(
        mountInfo, {
          appConfig: appConfig
        })
      // 区分单条挂载与批量挂载
      if (this.curIdList.length === 0) {
        mountInfo.replicaId = this.curBuildId
        this.$Modal.confirm({
          title: '确认挂载？',
          onOk: () => {
            this.$Notice.info({
              title: '请求中……',
              desc: '正在挂载，请等待'
            })
            const param = {
              mountInfo: mountInfo
            }
            mountBackup(param).then(({ data: { code, message } }) => {
              if (code !== 200) {
                this.$Notice.error({
                  title: '错误',
                  desc: message
                })
                return
              }
              this.$Notice.success({
                title: '成功',
                desc: '请求成功,请稍后查看挂载结果'
              })
              this.mountModal = false
              this.getData(this.currentPage)
            }).catch(err => {
              if (err && err.response) {
                this.$Notice.error({
                  title: '错误',
                  desc: err.response.data.msg
                })
              }
            })
          },
          onCancel: () => {
          }
        })
      } else {
        let mountInfoList = []
        for (let i in this.curIdList) {
          var tmpMap = Object.assign({}, mountInfo);
          tmpMap.replicaId = this.curIdList[i]
          mountInfoList.push(tmpMap)
        }
        this.$Modal.confirm({
          title: '确认挂载？',
          onOk: () => {
            this.$Notice.info({
              title: '请求中……',
              desc: '正在批量挂载，请等待'
            })
            const param = {
              mountInfos: mountInfoList
            }
            mountBackupBatch(param).then(({ data: { code, message } }) => {
              if (code !== 200) {
                this.$Notice.error({
                  title: '错误',
                  desc: message
                })
                return
              }
              this.$Notice.success({
                title: '成功',
                desc: '请求成功,请稍后查看挂载结果'
              })
              this.curIdList = []
              this.mountModal = false
              this.getData(this.currentPage)
            }).catch(err => {
              if (err && err.response) {
                this.$Notice.error({
                  title: '错误',
                  desc: err.response.data.msg
                })
              }
            })
          },
          onCancel: () => {
          }
        })
      }
    },
    doMountFc (name) { // 挂载FusionCompute平台
      let mountInfo = {
        replicaId: 0,
        targetType: Number(this.formMount.platform),
        targetId: Number(this.formMount.target),
        mountType: 1
      }
      /**
      if (!this.formMount.isCrtFcVm) {
        this.$Message.error('请选择是否创建虚拟机')
        return
      }
       **/
      Object.assign(
        mountInfo, {
          isExecScript: false
        })
      let isRefresh
      isRefresh = this.formMount.isRefresh === 'true'
      let isRegisterVM
      isRegisterVM = this.formMount.isCrtFcVm === 'true'
      let appConfig = {
        fcHostName: this.formMount.fcHostInfo.label,
        fcHostUrn: this.formMount.fcHostInfo.value,
        fcHostIp: this.formMount.fcHostInfo.ip,
        isRefresh: isRefresh,
        isRegisterVM: isRegisterVM,
        os: this.formMount.os
      }
      Object.assign(
        mountInfo, {
          appConfig: appConfig
        })
      // 区分单条挂载与批量挂载
      if (this.curIdList.length === 0) {
        mountInfo.replicaId = this.curBuildId
        this.$Modal.confirm({
          title: '确认挂载？',
          onOk: () => {
            this.$Notice.info({
              title: '请求中……',
              desc: '正在挂载，请等待'
            })
            const param = {
              mountInfo: mountInfo
            }
            mountBackup(param).then(({ data: { code, message } }) => {
              if (code !== 200) {
                this.$Notice.error({
                  title: '错误',
                  desc: message
                })
                return
              }
              this.$Notice.success({
                title: '成功',
                desc: '请求成功,请稍后查看挂载结果'
              })
              this.mountModal = false
              this.getData(this.currentPage)
            }).catch(err => {
              if (err && err.response) {
                this.$Notice.error({
                  title: '错误',
                  desc: err.response.data.msg
                })
              }
            })
          },
          onCancel: () => {
          }
        })
      } else {
        let mountInfoList = []
        for (let i in this.curIdList) {
          var tmpMap = Object.assign({}, mountInfo);
          tmpMap.replicaId = this.curIdList[i]
          mountInfoList.push(tmpMap)
        }
        this.$Modal.confirm({
          title: '确认挂载？',
          onOk: () => {
            this.$Notice.info({
              title: '请求中……',
              desc: '正在批量挂载，请等待'
            })
            const param = {
              mountInfos: mountInfoList
            }
            mountBackupBatch(param).then(({ data: { code, message } }) => {
              if (code !== 200) {
                this.$Notice.error({
                  title: '错误',
                  desc: message
                })
                return
              }
              this.$Notice.success({
                title: '成功',
                desc: '请求成功,请稍后查看挂载结果'
              })
              this.curIdList = []
              this.mountModal = false
              this.getData(this.currentPage)
            }).catch(err => {
              if (err && err.response) {
                this.$Notice.error({
                  title: '错误',
                  desc: err.response.data.msg
                })
              }
            })
          },
          onCancel: () => {
          }
        })
      }
    },
    doMountCus (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          let mountInfo = {
            targetType: Number(this.formMount.platform),
            targetId: Number(this.formMount.target),
            mountType: Number(this.mountTypeMap[this.mountTypeText]),
          }
          // if(this.formMount.mountPoint){
          //   Object.assign(
          //     mountInfo, {
          //       appconfig:{
          //         mountPoint:this.formMount.mountPoint
          //       }
          //     })
          // }
          if(this.mountTypeText==='应用默认挂载'){
            if(this.scriptName){
              Object.assign(
                mountInfo, {
                  isExecScript: true,
                  scriptId:Number(this.scriptName)
                })
            }else {
              Object.assign(
                mountInfo, {
                  isExecScript: false
                })
            }
          }else if(this.mountTypeText==='仅导出共享路径'){
            Object.assign(
              mountInfo, {
                isExecScript: false,
              })
          }else {
            this.$Message.error('挂载类型错误')
          }

          // 区分单条挂载与批量挂载
          if(this.curIdList.length===0){
            mountInfo.replicaId=this.curBuildId
            this.$Modal.confirm({
              title: '确认挂载？',
              onOk: () => {
                this.$Notice.info({
                  title: '请求中……',
                  desc: '正在挂载，请等待'
                })
                const param = {
                  mountInfo: mountInfo
                }
                mountBackup(param).then(({ data: { code, message } }) => {
                  if (code !== 200) {
                    this.$Notice.error({
                      title: '错误',
                      desc: message
                    })
                    return
                  }
                  this.$Notice.success({
                    title: '成功',
                    desc: '请求成功,请稍后查看挂载结果'
                  })
                  this.mountModal = false
                  this.getData(this.currentPage)
                }).catch(err => {
                  if (err && err.response) {
                    this.$Notice.error({
                      title: '错误',
                      desc: err.response.data.msg
                    })
                  }
                })
              },
              onCancel: () => {
              }
            })
          }else {
            let mountInfoList=[]
            for(let i in this.curIdList){
              var tmpMap = Object.assign({}, mountInfo);
              tmpMap.replicaId=this.curIdList[i]
              mountInfoList.push(tmpMap)
            }

            this.$Modal.confirm({
              title: '确认挂载？',
              onOk: () => {
                this.$Notice.info({
                  title: '请求中……',
                  desc: '正在批量挂载，请等待'
                })
                const param = {
                  mountInfos: mountInfoList
                }
                mountBackupBatch(param).then(({ data: { code, message } }) => {
                  if (code !== 200) {
                    this.$Notice.error({
                      title: '错误',
                      desc: message
                    })
                    return
                  }
                  this.$Notice.success({
                    title: '成功',
                    desc: '请求成功,请稍后查看挂载结果'
                  })
                  this.curIdList=[]
                  this.mountModal = false
                  this.getData(this.currentPage)
                }).catch(err => {
                  if (err && err.response) {
                    this.$Notice.error({
                      title: '错误',
                      desc: err.response.data.msg
                    })
                  }
                })
              },
              onCancel: () => {
              }
            })
          }
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    unMount (row) {
      this.$Notice.info({
        title: '请求中……',
        desc: '正在卸载，请等待'
      })
      const param = {
        id: row.id
      }
      unmountBackup(param).then(({ data: { code, message } }) => {
        if (code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: message
          })
          return
        }
        this.$Notice.success({
          title: '成功',
          desc: '请求成功,请稍后查看卸载结果'
        })
        this.mountModal = false
        this.getData(this.currentPage)
      }).catch(err => {
        if (err && err.response) {
          // console.error(err)
          this.$Notice.error({
            title: '错误',
            desc: err.response.data.msg
          })
        }
      })
    },
    openMirror (row) {
      this.formMirror.backupId = row.id
      this.mirrorModal = true
    },
    crtMirror (name) {
      this.$refs[name].validate((valid) => {
        if (valid) {
          this.isdisabledFn = true
          const param = {
            id: this.formMirror.backupId,
            name: this.formMirror.name,
            desc: this.formMirror.desc
          }
          this.$Notice.info({
            title: '请求中……',
            desc: '正在生成镜像，请等待',
            duration: 10
          })
          crtMirror(param).then(({ data: { code, message } }) => {
            if (code !== 200) {
              this.$Notice.error({
                title: '错误',
                desc: message
              })
              this.isdisabledFn = false
              return
            }
            this.$Notice.success({
              title: '成功',
              desc: '请求成功'
            })
            this.mirrorModal = false
            this.isdisabledFn = false
            this.getData()
          }).catch(err => {
            if (err && err.response) {
              // console.error(err)
              this.$Notice.error({
                title: '错误',
                desc: err.response.data.msg
              })
              this.isdisabledFn = false
            }
          })
        } else {
          this.$Message.error('请输入必填项!')
        }
      })
    },
    destroyBackup (row) {
      const params = Object.assign({}, {
        id: row.id
      })
      this.showSpin = true
      httpRequest(this, deleteBackup, [params]).then(res => {
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
        this.getData(this.currentPage)
      })
    },

    // 脚本列表
    getScriptList () {
      this.scriptList = []
      httpRequest(this, getAllScriptList, []).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        for (let i in res.response.data) {
          this.scriptList.push({ key: res.response.data[i].name, value: res.response.data[i].id })
        }
      })
    },

    getCasHostList () {
      this.casHostList = []
      const params = Object.assign({}, {
        id: Number(this.formMount.target)
      })
      httpRequest(this, getCasHosts, [params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if (Array.isArray(res.response.host)) {
          for (let i in res.response.host) {
            this.casHostList.push({ key: res.response.host[i].name, value: res.response.host[i].id })
          }
        } else {
          this.casHostList.push({ key: res.response.host.name, value: res.response.host.id })
        }
      })
    },
    getFcHostList () {
      this.fcHostList = []
      const params = Object.assign({}, {
        id: Number(this.formMount.target)
      })
      httpRequest(this, getFcHosts, [params]).then(res => {
        if (res.code !== 200) {
          this.$Notice.error({
            title: '错误',
            desc: res.message
          })
          return
        }
        if (Array.isArray(res.response.hosts)) {
          for (let i in res.response.hosts) {
            if (res.response.hosts[i].status === 'normal') {
              this.fcHostList.push({
                label: res.response.hosts[i].name,
                value: res.response.hosts[i].urn,
                ip: res.response.hosts[i].ip
              })
            }
          }
        }
      })
    },
    get_replica_info () {
      let replicaId = this.curBuildId
      httpRequest(this, backupDetails, [replicaId]).then(res => {
        console.log('获取副本详情:', res.response)
        this.casDiskLis = res.response
      })
    },
    selectHost (value) {
      console.log(value)
    },
    backupDetail (row) {
      this.showDetail = true
      this.curBuildId = row.id
      this.$refs.backupInfo.detail(row.id)
      this.$refs.backupInfo.getRecordData(1, row.id)

    },
    openSnapshot (row) {
      this.showCreateSnapshot = true
      this.currentBackup = row.name
      this.currentId = row.id
    },
    cancel () {
      this.curIdList = []
    },
    // 搜索
    search () {
      this.currentPage = 1
      this.getData(this.currentPage, this.pageSize)
    },
    // 翻页
    changePage (page) {
      this.currentPage = page
      this.getData(page, this.pageSize)
    },
    changeSize (size) {
      this.pageSize = size
      this.getData(1, size)
    },
    closeCreateModal () {
      this.showCreateSnapshot = false
    },
    refresh () {
      this.getData(this.currentPage)
    }
  }
}
</script>

<style lang="less">
@import "../../view/env/css/ticket.less";

.create-step {
  padding: 20px 0;
  border-bottom: 2px solid #0099CC;
  margin-bottom: 20px;
}

/* 让方形复选框变成圆形单选框样式，最好在树组件外套一个父盒子，在样式前加父级选择器，以免影响其他树组件*/
/*.ivu-checkbox-inner {*/
/*  border-radius: 50%;*/
/*}*/
</style>
