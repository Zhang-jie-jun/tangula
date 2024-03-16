
<template>
  <div>
    <Col span="12">
      <Card title="基本信息" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <Row style="margin-bottom: 20px">
          <span class="text-class">fsid:  </span> <span class="text-class2">{{fsid}}</span>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">health: </span> <span v-if="health==='HEALTH_OK'" class="health-ok-class">{{health}} </span><span v-else class="health-warn-class">{{health}} </span>
          <Button v-if="health!=='HEALTH_OK'" size="small"style="margin-left: 30px" @click="openHealth">查看</Button>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">election_epoch: </span> <span class="text-class2">{{election_epoch}} </span>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">quorum_names: </span> <span class="text-class2">{{quorum_names}} </span>
        </Row>
      </Card>
    </Col>
    <Col span="12">
      <Card title="pgmap" icon="ios-send" style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <div style="width:400px;height: 400px;" id="pgmap_charts" ref="pgmapCharts" ></div>
      </Card>
    </Col>
<!--    <Col span="12">-->
<!--      <Card title="osdmap" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 300px;">-->
<!--        <Row style="margin-bottom: 20px">-->
<!--          <span class="text-class">num_osds: </span> <span class="text-class2">{{osdmap.num_osds}} </span>-->
<!--        </Row>-->
<!--        <Row style="margin-bottom: 20px">-->
<!--          <span class="text-class">num_up_osds: </span> <span class="text-class2">{{osdmap.num_up_osds}} </span>-->
<!--        </Row>-->
<!--        <Row style="margin-bottom: 20px">-->
<!--          <span class="text-class">num_in_osds: </span> <span class="text-class2">{{osdmap.num_in_osds}} </span>-->
<!--        </Row>-->
<!--        <Row style="margin-bottom: 20px">-->
<!--          <span class="text-class">num_remapped_pgs: </span> <span class="text-class2">{{osdmap.num_remapped_pgs}} </span>-->
<!--        </Row>-->
<!--      </Card>-->
<!--    </Col>-->
    <Col span="12">
      <Card title="mgrmap" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <Row style="margin-bottom: 20px">
          <span class="text-class">active_name: </span> <span class="text-class2"><Button type="success" style="margin-left: 30px;font-size: 16px" @click="openAvailable">{{mgrmap.active_name}}</Button> </span>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">active_addr: </span> <span class="text-class2">{{mgrmap.active_addr}} </span>
        </Row>
        <Row style="margin-bottom: 20px">
          <Col span="4"><span class="text-class">standbys: </span></Col>
          <Col span="4"><Table style="margin-bottom: 20px;margin-top:20px" :width="100" :show-header=false :columns="standbysTitle" :data="mgrmap.standbys" size="small"></Table></Col>
        </Row>
      </Card>
    </Col>
    <Col span="12">
      <Card title="monmap" icon="ios-send"  style=" margin-bottom: 20px;margin-right: 20px;height: 500px;">
        <Row style="margin-bottom: 20px">
          <span class="text-class">created: </span> <span class="text-class2">{{monmap.created}} </span>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">modified: </span> <span class="text-class2">{{monmap.modified}} </span>
        </Row>
        <Row style="margin-bottom: 20px">
          <span class="text-class">mons: </span>
        </Row>
        <Table :row-class-name="rowClassName" searchable border :columns="columnsTitle" :data="mons" size="small">
        </Table>
      </Card>
    </Col>
    <Modal
      v-model="healModal"
      title="health">
      <p>{{health_checks}}</p>
      <div slot="footer">
        <Button  @click="closeHealthModal">关闭</Button>
      </div>
    </Modal>
    <Modal
      v-model="availableModal"
      title="available_modules">
      <Table style="margin-bottom: 20px;" :row-class-name="rowClassName" :show-header=false border :columns="availableTitle" :data="availableData" size="small">
      </Table>
      <div slot="footer">
        <Button  @click="closeavailableModal">关闭</Button>
      </div>
    </Modal>

    <Modal
      v-model="standbyModal"
      title="available_modules">
      <Table style="margin-bottom: 20px" :row-class-name="rowClassName" border :columns="availableTitle" :data="standbysData" size="small">
      </Table>
      <div slot="footer">
        <Button  @click="closestandbyModal">关闭</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { httpRequest } from '@/libs/tools.js'
import tdTheme from '_c/dashboard/theme.json'
import { on, off } from '@/libs/tools'
import {getCephInfo, getAllJson, jsonToTree} from "@/api/env";
import { toGB } from '@/api/backup.js'
echarts.registerTheme('tdTheme', tdTheme)
export default {
  name: 'home',
  data () {
    return {
      resData: [],
      healModal:false,
      availableModal:false,
      standbyModal:false,
      cephData:[],
      mons:[],
      availableData:[],
      standbysData:[],
      fsid:"",
      health:"",
      health_checks:"",
      election_epoch:"",
      quorum_names:"",
      monmap:{
        epoch:"",
        created:"",
        modified:"",

      },
      osdmap:{
        epoch: 0,
        num_osds: 0,
        num_up_osds: 0,
        num_in_osds: 0,
        num_remapped_pgs: 0
      },
      mgrmap:{
        epoch: "",
        active_gid: "",
        active_name: "",
        active_addr: "",
        standbys:[]

      },
      pgmap:{
        num_pgs:"",
        num_pools:"",
        bytes_total:""
      },
      columnsTitle: [
        {
          title: 'rank',
          width: 70,
          key: 'rank',
          align: 'center'
        },
        {
          title: 'name',
          key: 'name',
          align: 'center'
        },
        {
          title: 'addr',
          key: 'addr',
          align: 'center'
        },
        {
          title: 'public_addr',
          key: 'public_addr',
          align: 'center'
        }
      ],
      availableTitle: [
        {
          title: 'name',
          key: 'name',
          align: 'center'
        }
      ],
      standbysTitle: [
        {
          title: 'name',
          key: '',
          align: 'center',
          style: 'padding: 0px',
          render: (h, params) => {
            let atts = []
            atts.push(
              h('Button', {
                props: {
                  size: 'small',
                  type: 'warning'
                },
                style: {
                  marginLeft: '1%'
                },
                on: {
                  click: () => {
                    this.openStandby(params.row.available_modules)
                  }
                }
              }, params.row.name)
            )
            return h('div', atts)
          }
        }
      ],
      option: {
        title: {
          text: '',
          left: 'center'
        },
        tooltip: {
          trigger: 'item'
        },
        legend: {
          top: '5%',
          left: 'center'
        },
        series: [
          {
            name: '大小(GB)',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '20',
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: []
          }
        ]
      }
    }
  },
  methods: {
    updateDom () {
      let dbCount = echarts.init(document.getElementById('pgmap_charts'))
      dbCount.setOption(this.option)
    },
    getData () {
      httpRequest(this, getCephInfo, []).then(res => {
        if (res.code !== 200) {
          this.$Message.error('获取数据有误')
          return
        }

        this.fsid=res.response.fsid
        this.health=res.response.health.status
        this.health_checks=JSON.stringify(res.response.health.checks)
        this.election_epoch=res.response.election_epoch
        let quorum=""
        for(let i in res.response.quorum_names){
          quorum+=res.response.quorum_names[i]
          if(i<res.response.quorum_names.length-1){
            quorum+="，"
          }
        }
        this.quorum_names=quorum
        this.mons=res.response.monmap.mons

        this.osdmap.epoch=res.response.osdmap.epoch
        this.osdmap.num_osds=res.response.osdmap.num_osds
        this.osdmap.num_up_osds=res.response.osdmap.num_up_osds
        this.osdmap.num_in_osds=res.response.osdmap.num_in_osds
        this.osdmap.num_remapped_pgs=res.response.osdmap.num_remapped_pgs

        this.monmap.epoch=res.response.monmap.epoch
        this.monmap.created=res.response.monmap.created
        this.monmap.modified=res.response.monmap.modified

        this.mgrmap.epoch=res.response.mgrmap.epoch
        this.mgrmap.active_gid=res.response.mgrmap.active_gid
        this.mgrmap.active_name=res.response.mgrmap.active_name
        this.mgrmap.active_addr=res.response.mgrmap.active_addr
        this.mgrmap.standbys=res.response.mgrmap.standbys
        let available_names=[]
        for(let i in res.response.mgrmap.available_modules){
          available_names.push({
            name:res.response.mgrmap.available_modules[i]
          })
        }
        this.availableData=available_names

        this.option.title.text='总容量:'+toGB(res.response.pgmap.bytes_total)
        this.option.series[0].data=[
          {
            value:toGB(res.response.pgmap.bytes_used),
            name:'已使用'
          },
          {
            value:toGB(res.response.pgmap.bytes_avail),
            name:'未使用'
          }
        ]

        this.updateDom()
      })
    },
    openHealth(){
      this.healModal=true
    },
    openAvailable(){
      this.availableModal=true
    },
    openStandby(available_modules){
      this.standbysData=available_modules
      this.standbyModal=true
    },
    rowClassName (row, index) {
        return 'demo-table-info-row'
    },
    closeHealthModal(){
      this.healModal = false
    },
    closeavailableModal(){
      this.availableModal = false
    },
    closestandbyModal(){
      this.standbyModal = false
    },

  },
  mounted: function () {
    this.getData()
    this.$nextTick(() => {
      this.dom = echarts.init(this.$refs.pgmapCharts)
      this.dom.setOption(this.option)
      on(window, 'resize', this.resize)
    })
  },
  beforeDestroy () {
    off(window, 'resize', this.resize)
  }
}
</script>

<style lang="less">
.count-style{
  margin-bottom: 20px;
  margin-right: 20px;
  height: 400px;
}
.text-class{
  font-size: 18px;
  margin-bottom: 20px;
  color: #0099CC;
}
.text-class2{
  font-size: 18px;
  margin-bottom: 20px;
  color: #000000;
}
.health-ok-class{
  font-size: 18px;
  margin-bottom: 20px;
  color: #339933;
}
.health-warn-class{
  font-size: 18px;
  margin-bottom: 20px;
  color: #FF9900;
}
.ivu-table .demo-table-info-row td{
  background-color: #99CCFF;
  color: #fff;
  font-size: 16px;
}
.ivu-table .demo-table-error-row td{
  background-color: #ff6600;
  color: #fff;
}
</style>
