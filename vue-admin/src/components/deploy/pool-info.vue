<template>
  <div>
    <Collapse v-model="showPanel">
      <Panel name="build">
        存储池信息
        <p slot="content">
        <Form :label-width="90" label-position="left" class="not-validate-item">
          <Row>
            <Col span="12">
              <FormItem label="存储池名称">
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
import {getPoolDetails} from '@/api/env.js'
import { httpRequest } from '@/libs/tools.js'
export default {
  name: 'poolInfo',
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
    }
  },
  computed: {
    deployColor () {
      return '#515a6e'
    }
  },
  methods: {
    detail (_id) {
      this.showDetail = true
      httpRequest(this, getPoolDetails, [_id]).then(res => {
        this.detailData = res.response
      })
    }
  }
}
</script>
