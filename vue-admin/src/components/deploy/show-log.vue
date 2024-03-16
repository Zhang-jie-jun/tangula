<template>
  <div>
    <!-- @on-cancel="closeModal" -->
    <Modal v-model="show" width="1024" :draggable="true" class="log-ivu-modal-body log-ivu-modal-header" footer-hide>
      <p slot="header" style="min-height:32px;">
        <Col span='12'>
          <span style="margin-right: 20px">日志</span>
        </Col>
<!--        <Col span='4'>-->
<!--          <Input v-model="logKeyword" placeholder="Enter something..." style="width: 200px" />-->
<!--        </Col>-->
        <Col span='10'>
          <Button type="success" size="small" icon="md-refresh" shape="circle" @click="refresh" style="float:right"></Button>
<!--          <Button loading shape="circle" type="primary"  @click="refresh" style="float:right"></Button>-->
        </Col>
      </p>
      <highlight-text id="showContentCom" :data="logData" :keyword="logKeyword" :overWriteStyle="{color: logKeywordColor}"></highlight-text>
    </Modal>
  </div>
</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import HighlightText from '_c/high-light'
import {getInstanceDetails} from "@/api/backup";

export default {
  name: 'DeployLog',
  components: {
    HighlightText
  },
  data () {
    return {
      currentId: "",
      env: '',
      show: false,
      logData: '',
      logKeyword: '',
      logKeywordColor: '#339933',
      showBuildDetail: false,
      kibanaUrl: ''

    }
  },
  methods: {
    closeModal () {
      // this.show = false
    },
    deployLog (logId) {
      this.logData = ''
      this.currentId=logId
      const params = Object.assign({}, {
        id: logId
      })
      httpRequest(this, getInstanceDetails, [params]).then(res => {
        if(res.response.data){
          for(let i=0;i<res.response.data.length;i++){
            this.logData += res.response.data[i].created_time + '=>' + res.response.data[i].info
            if(res.response.data[i].detail !== ''){
              this.logData += '：'+ res.response.data[i].detail
            }
            this.logData += '\r\n'
          }
        }
      })
    },
    setLogParams (state = true, logId) {
      this.logData = ''
      this.deployLog(logId)
      this.show = state
    },
    refresh(){
      this.deployLog(this.currentId);
    }
  },
  mounted: function () {
  }
}
</script>
