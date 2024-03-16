<template>
  <div>
    <Collapse v-model="showPanel">
      <Panel name="build">
        平台信息
        <p slot="content">
        <Form :label-width="120" label-position="left" class="not-validate-item">
          <Row>
            <Col span="12">
              <FormItem label="平台名">
                <span>{{ detailData.name }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="描述">
                <span>{{ detailData.desc }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="类型">
                <span>{{ typeText }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="版本">
                <span>{{ detailData.version }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="ip">
                <span>{{ detailData.ip }}</span>
              </FormItem>
            </Col>
            <Col span="12">
              <FormItem label="端口">
                <span>{{ detailData.port }}</span>
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
              <FormItem label="用户名">
                <span>{{ detailData.username }}</span>
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
import { getPlatDetails } from '@/api/env.js'
import { httpRequest } from '@/libs/tools.js'
export default {
  name: 'platformInfo',
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
      typeText: ''
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
      httpRequest(this, getPlatDetails, [_id]).then(res => {
        this.detailData = res.response
        this.typeText = this.typeMap[res.response.type]
        // this.deployDetail = res.deploy
      })
    }
  }
}
</script>
