<template>
  <Row :gutter="20">
    <Form :label-width="100" label-position="right" class="per-form-width">
      <FormItem label="用户名">
        <Input v-model="userName" readonly="readonly"/>
      </FormItem>
      <FormItem label="用户角色">
        <Input v-model="roleChn" readonly="readonly"/>
      </FormItem>
      <FormItem label="业务线">
         <Input v-model="business" readonly="readonly"/>
        <!--<Select v-model="project" style="width:200px" @on-change="selectProject()">-->
          <!--<Option v-for="item in projects" :value="item.id" :key="item.label">{{ item.label }}</Option>-->
        <!--</Select>-->
      </FormItem>
    </Form>
    <Row class="per-table-margin">
      <!--<Table border :columns="columnsTitle" :data="apps" stripe size="small"></Table>-->
      <!--<div class="per-page">-->
        <!--<div class="per-page-position">-->
            <!--<Page :total="total" :current="1" show-total show-elevator @on-change="changePage"></Page>-->
        <!--</div>-->
      <!--</div>-->
    </Row>
  </Row>
</template>
<script>

import { getApplication } from '@/api/permission'
import { httpRequest } from '@/libs/tools.js'
import { mapState } from 'vuex'

export default {
  name: 'PermInfo',
  // 通过父组件传过来的数据，
  // 子组件只能读
  props: {
    recodeId: {
      type: [String, Number],
      default: ''
    }
  },
  // 子组件自己的数据定义，可读可写
  data () {
    return {
      total: 0,
      proIds: [],
      apps: [],
      project: '',
      proStr: '',
      columnsTitle: [
        // {
        //   type: 'selection',
        //   width: 60,
        //   align: 'center'
        // },
        {
          title: '项目',
          key: 'project_chn',
          align: 'center'
        },
        {
          title: '应用名',
          key: 'app_name',
          align: 'center'
        },
        {
          title: '应用中文名',
          key: 'app_name_chn',
          align: 'center'
        }
      ]
    }
  },
  computed: {
    ...mapState({
      projects: state => state.perm.projects,
      roleChn: state => state.perm.roleChn,
      userName: state => state.perm.userName,
      business: state => state.perm.business
    })
  },
  methods: {
    getApps (page = 1, project_id = null) {
      this.proIds = []
      if (this.projects.length === 0) return
      if (this.project !== '') project_id = this.project

      if (project_id === null) {
        this.projects.forEach(ele => {
          this.proIds.push(ele.id)
        })
      } else {
        this.proIds = [project_id]
      }

      let pros = this.proIds.join(',')
      httpRequest(this, getApplication, ['', pros, page]).then(res => {
        this.apps = res.results
        this.total = res.count
      })
    },
    selectProject () {
      this.getApps(1, this.project)
    },
    // 翻页
    changePage (page) {
      this.getApps(page)
    }
  },

  mounted: function () {
    // console.log(this.proIds, '---------')
    // console.log(this.roleChn)
    // this.proStr = this.projectsChn.join('/')
    this.getApps()
  }

}
</script>

<style lang="less">
.per-card-min-height {
  min-height: 580px;
}
.per-form-width {
  width: 300px;
}
.per-table-margin {
  padding-left: 16px;
  padding-right: 16px;
}
.per-page {
  margin-top: 10px;
  margin-bottom: 5px;
}
.per-page-position {
  float: right;
}
</style>
