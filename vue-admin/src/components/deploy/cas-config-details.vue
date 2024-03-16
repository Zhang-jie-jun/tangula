<template>
  <div>
    <!-- @on-cancel="closeModal" -->
    <Modal v-model="show"  width="1024" :draggable="true" class="log-ivu-modal-body log-ivu-modal-header" footer-hide>
      <p slot="header" style="min-height:32px;">
        <Icon type="md-information-circle"></Icon>
        <Col span='17'>
          <span style="margin-right: 20px">脚本</span>
        </Col>
        <Col span='4'>
          <Input v-model="logKeyword" placeholder="Enter something..." style="width: 200px" />
        </Col>
      </p>
      <highlight-text id="showContentCom" :data="fileData" :keyword="logKeyword" :overWriteStyle="{color: logKeywordColor}"></highlight-text>
    </Modal>
  </div>
</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import HighlightText from '_c/high-light'
import {getCasFileDetails} from "@/api/backup";

export default {
  name: 'casConfigDetails',
  components: {
    HighlightText
  },
  data () {
    return {
      deployId: -1,
      env: '',
      show: false,
      fileData: '',
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
    fileLog (logId) {
      const params = Object.assign({}, {
        id: logId
      })
      httpRequest(this, getCasFileDetails, [params]).then(res => {
        if(res){
          this.fileData=res
        }
      })
    },
    setFileParams (state = true, logId) {
      this.fileData = ''
      this.fileLog(logId)
      this.show = state
    }
  },
  mounted: function () {
  }
}
</script>
