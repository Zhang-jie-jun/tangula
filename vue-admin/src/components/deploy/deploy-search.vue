<template>
  <div>
    <Row class="search-margin-botton">
      <!--<Col span="3" class="col-margin">-->
        <!--<Select v-model="env" clearable @on-change="conditionChange" placeholder="请选择环境">-->
          <!--<Option v-for="(value, key) in searchAppEnvs" :value="key" :key="key" >{{ value }}</Option>-->
        <!--</Select>-->
      <!--</Col>-->
      <Col span="3" class="col-margin">
        <Select v-model="status" clearable @on-change="conditionChange" placeholder="请选择状态">
          <Option v-for="(value, key) in searchAppStates" :value="key" :key="value" >{{ value }}</Option>
        </Select>
      </Col>
      <Col span="3" class="col-margin">
        <Input v-model="appName" enter-button @on-enter="conditionChange" placeholder="应用名..."/>
      </Col>
      <Col span="4" class="col-margin">
        <DatePicker v-model="startTime" type="datetime" format="yyyy-MM-dd HH:mm:ss" @on-ok="conditionChange" placeholder="开始时间"></DatePicker>
      </Col>
      <Col span="4" class="col-margin">
        <DatePicker v-model="endTime" type="datetime" format="yyyy-MM-dd HH:mm:ss" @on-ok="conditionChange" placeholder="结束时间"></DatePicker>
      </Col>
      <Button class="right-but" type="warning" @click="exportData">导出</Button>
      <Button class="right-but" type="primary" @click="conditionChange">搜索</Button>
    </Row>
  </div>
</template>

<script>
import { formatDate } from '@/libs/tools.js'
export default {
  name: 'DeploySearch',
  props: {
  },
  data () {
    return {
      appName: '',
      env: '',
      status: '',
      startTime: '',
      endTime: '',
      searchAppEnvs: this.$config.deployEnv,
      searchAppStates: this.$config.deployState
    }
  },
  methods: {
    clickExport () {
      this.$emit('export-click')
    },
    conditionChange () {
      if (this.startTime !== '' && this.endTime !== '' && new Date(this.startTime) >= new Date(this.endTime)) {
        this.$Message.warning('开始时间不能大于等于结束时间！')
        return
      }

      let params = {
        app_name: this.appName,
        env: this.env,
        status: this.status,
        start_time: this.startTime,
        end_time: this.endTime
      }

      if (this.startTime) params.start_time = formatDate(this.startTime, 'yyyy-MM-dd hh:mm:ss')
      if (this.endTime) params.end_time = formatDate(this.endTime, 'yyyy-MM-dd hh:mm:ss')

      this.$emit('deploy-search', params)
    },
    exportData () {
      this.$emit('deploy-search-export')
    }
  }
}
</script>
