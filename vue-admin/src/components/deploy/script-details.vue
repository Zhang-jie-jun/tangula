<template>
  <div>
    <!-- @on-cancel="closeModal" -->
    <Modal v-model="show" width="1024" :draggable="true" class="log-ivu-modal-body log-ivu-modal-header" footer-hide>
      <p slot="header" style="min-height:32px;">
        <Icon type="md-information-circle"></Icon>
        <Col span='17'>
          <span style="margin-right: 20px">脚本</span>
        </Col>
        <Col span='4'>
          <Input v-model="logKeyword" placeholder="Enter something..." style="width: 200px" />
        </Col>
      </p>
      <highlight-text id="showContentCom" :data="logData" :keyword="logKeyword" :overWriteStyle="{color: logKeywordColor}"></highlight-text>
    </Modal>
  </div>
</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import HighlightText from '_c/high-light'
import {getScriptDetails} from "@/api/backup";

export default {
  name: 'ScriptDetails',
  components: {
    HighlightText
  },
  data () {
    return {
      deployId: -1,
      env: '',
      show: false,
      logData: '',
      logKeyword: '',
      logKeywordColor: '#1972FE',
      showBuildDetail: false,
      kibanaUrl: ''

    }
  },
  methods: {
    closeModal () {
      // this.show = false
    },
    deployLog (logId) {
      const params = Object.assign({}, {
        id: logId
      })
      httpRequest(this, getScriptDetails, [params]).then(res => {
        if(res){
          this.logData=res
        }
      })
    },
    setLogParams (state = true, logId) {
      this.logData = ''
      this.deployLog(logId)
      this.show = state
    }
  },
  mounted: function () {
  }
}
</script>
