<template>
  <div>
    <Collapse v-model="showPanel">
      <Panel name="build">
        镜像信息
        <p slot="content">
        <Form :label-width="120" label-position="left" class="not-validate-item">
          <Row>
            <Col span="12">
              <FormItem label="镜像名称">
                <span>{{ detailData.name }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="描述">
                <span>{{ detailData.desc }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="uuid">
                <span>{{ detailData.uuid }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="类型">
                <span>{{ typeText }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="大小">
                <span>{{ detailData.size }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="状态">
                <span>{{ statusText }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="所属池">
                <span>{{ poolName }}</span>
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
        </Form>
        </p>
      </Panel>
    </Collapse>
  </div>
</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import {getImageStatus, mirrorDetails} from "@/api/backup";
export default {
  name: 'mirrorInfo',
  props: {
    bId: {
      type: [Number, String],
      required: false
    },
    typeMap: Object
  },
  data () {
    return {
      detailData: {},
      deployDetail: {},
      showDetail: false,
      showPanel: 'build',
      typeText: '',
      statusText:'',
      poolName:'',
    }
  },
  computed: {
    deployColor () {
      return '#515a6e'
      // if (!this.detailData.deploy) return '#515a6e'
      // let color = choiceDeployStateColor(this.detailData.deploy.status_text)
      // return color
    }
  },
  methods: {
    detail (_id) {
      this.showDetail = true
      httpRequest(this, mirrorDetails, [_id]).then(res => {
        this.detailData = res.response
        this.typeText = this.typeMap[res.response.type]
        this.poolName = res.response.pool.name
        this.statusText = getImageStatus(res.response.status)
        // this.deployDetail = res.deploy
      })
    }
  }
}
</script>
