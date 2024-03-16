<template>
  <div>
    <Col span="12">
      <Card title="平台信息" icon="ios-send" style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <div style="height: 400px;" id="resource_charts" ref="resourceCharts" ></div>
      </Card>
    </Col>
    <Col span="12">
      <Card title="用户信息" icon="ios-send" style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <div style="height: 400px;" id="user_charts" ref="userCharts" ></div>
      </Card>
    </Col>
    <Col span="12">
      <Card title="操作系统" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 300px;">
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">os:</span>
          </Col>
          <Col span="12">
            <span class="txt-class2">{{os.goos}}</span>
          </Col>
        </Row>
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">numCpu:</span>
          </Col>
          <Col span="12">
            <span class="txt-class2">{{os.numCpu}}</span>
          </Col>
        </Row>
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">compiler:</span>
          </Col>
          <Col span="12">
            <span class="txt-class2">{{os.compiler}}</span>
          </Col>
        </Row>
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">goVersion:</span>
          </Col>
          <Col>
            <span class="txt-class2">{{os.goVersion}}</span>
          </Col>
        </Row>
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">numGoroutine:</span>
          </Col>
          <Col>
            <span class="txt-class2">{{os.numGoroutine}}</span>
          </Col>
        </Row>
      </Card>
    </Col>
    <Col span="12">
      <Card title="硬盘" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 300px;">
        <Col span="12">
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">total (MB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{disk.totalMb}}</span>
            </Col>
          </Row>
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">used (MB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{disk.usedMb}}</span>
            </Col>
          </Row>
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">total (GB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{disk.totalGb}}</span>
            </Col>
          </Row>
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">used (GB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{disk.usedGb}}</span>
            </Col>
          </Row>
        </Col>
        <Col span="12">
          <i-circle :percent="disk.usedPercent">
            <span class="demo-Circle-inner" style="font-size:24px">{{disk.usedPercent}}%</span>
          </i-circle>
        </Col>
      </Card>
    </Col>
    <Col span="12">
      <Card title="CPU" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 320px;">
        <Row style="margin-bottom: 20px">
          <Col span="12">
            <span class="txt-class">核心数量:</span>
          </Col>
          <Col span="12">
            <span class="txt-class2">{{cpu.cores}}</span>
          </Col>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="txt-class">运行中:</span>
        </Row>
        <Row>
          <Col v-if="core!==0" span="12"
               v-for="(core,index) in cpuUsedList"
               :key="index"
               :core="core"
          >
            <Progress :percent="core" status="active" />
          </Col>
        </Row>

      </Card>
    </Col>
    <Col span="12">
      <Card title="内存" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 320px;">
        <Col span="12">
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">total (MB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{ram.totalMb}}</span>
            </Col>
          </Row>
          <Row style="margin-bottom: 20px">
            <Col span="12">
              <span class="txt-class">used (MB):</span>
            </Col>
            <Col span="12">
              <span class="txt-class2">{{ram.usedMb}}</span>
            </Col>
          </Row>
        </Col>
        <Col span="12">
          <i-circle :percent="ram.usedPercent">
            <span class="demo-Circle-inner" style="font-size:24px">{{ram.usedPercent}}%</span>
          </i-circle>
        </Col>
      </Card>
    </Col>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { mapState } from 'vuex'
import {httpRequest, on} from "@/libs/tools";
import {getDashboardInfo, getSystemInfo} from "@/api/env";
import {toGB} from "@/api/backup";
import {createSocket, sendWSPush} from "@/libs/new-socket";
import JFSocketIo from "@/libs/jf-socket-io";
export default {
  name: 'home',
  data () {
    return {
      os: {},
      disk:{},
      ram:{},
      cpu:{},
      cpuUsedList:[],
      todayUser:0,
      option:{
        legend: {},
        tooltip: {},
        color:['#339933','#006699'],
        dataset: {
          source: []
        },
        xAxis: { type: 'category' },
        yAxis: {},
        // Declare several bar series, each will be mapped
        // to a column of dataset.source by default.
        series: [
          { type: 'bar' },
          { type: 'bar' }
        ]
      },
      option1:{
        color:['#ee6666','#5470c6'],
        title: {
          text: '',
          left: 'center'
        },
        tooltip: {
          trigger: 'item'
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: '用户',
            type: 'pie',
            radius: '50%',
            data: [],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      }
    }
  },
  computed: {
    ...mapState({
      userName: state => state.user.userName,
      defBusi: state => state.user.userBusi,
      permits: state => state.user.permits
    })
  },
  methods: {
    getSysData () {
      httpRequest(this, getSystemInfo, []).then(res => {
        this.os=res.response.os
        this.disk=res.response.disk
        this.ram=res.response.ram
        this.cpu=res.response.cpu
        let tmpList=[]
        for(let i in res.response.cpu.cpus){
          tmpList.push(Math.round(res.response.cpu.cpus[i]))
        }
        this.cpuUsedList=tmpList

      })
    },
    getPlatData () {
      httpRequest(this, getDashboardInfo, []).then(res => {
        let platformList=['平台',res.response.platformStat.public,res.response.platformStat.private]
        let poolList=['存储池',res.response.poolStat]
        let hostList=['主机',res.response.hostStat.public,res.response.hostStat.private]
        let replicaList=['副本',0.0,res.response.replicaStat]
        let imageList=['镜像',res.response.imageStat.public,res.response.imageStat.private]
        let scriptList=['脚本',res.response.scriptStat]
        this.option.dataset.source= [['类型','公共','私有'],poolList,platformList,hostList,replicaList,imageList,scriptList]

        this.option1.title.text='今日访问人数:'+res.response.userStat.todayNum
        let userData=[]
        userData.push(
          {
            value:res.response.userStat.activeNum,
            name:'活跃用户'
          },
          {
            value:res.response.userStat.totalNum,
            name:'总用户'
          }
        )
        this.option1.series[0].data=userData
        try {
          this.updateDom()
        }catch (err){

        }

      })
    },
    updateDom () {
      let upd = echarts.init(document.getElementById('resource_charts'))
      upd.setOption(this.option)

      let upd1 = echarts.init(document.getElementById('user_charts'))
      upd1.setOption(this.option1)
    },
    getWebsocket(){
      let newSocket=new JFSocketIo('home')
      newSocket.sendWS('home websockets')

      // 接收消息
      const getsocketData = e => {  // 创建接收消息函数
        const data = e && e.detail.data
        if(String(data)==='heart') {
          //console.log('websocket推送', data)
          this.getSysData()
          this.getPlatData()
        }
      }
      // 注册监听事件
      window.addEventListener('onmessageWS', getsocketData)

    }
  },
  mounted () {
    this.getSysData()
    this.getPlatData()
    this.$nextTick(() => {
      this.dom = echarts.init(this.$refs.resourceCharts)
      this.dom.setOption(this.option)

      this.dom = echarts.init(this.$refs.userCharts)
      this.dom.setOption(this.option1)
      on(window, 'resize', this.resize)
    })
    //this.getWebsocket()
  }
}
</script>

<style lang="less">
.txt-class{
  font-size: 18px;
  margin-bottom: 20px;
  color: #0099CC;
}
.txt-class2{
  font-size: 18px;
  margin-bottom: 20px;
  color: #000000;
}

</style>
