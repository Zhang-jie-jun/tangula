<template>
  <Row>
    <Col span="24" style="margin-bottom: 10px">
      <Card>
        <p slot="title">基本信息</p>
        <p><span class='demo-split-pane'>应用名：</span><span class="demo-split">{{ detailData.app_name }}</span></p>
        <p><span class='demo-split-pane'>版本号：</span><span class="demo-split">{{ detailData.app_version }}</span></p>
        <p><span class='demo-split-pane'>分支：</span><span class="demo-split">{{ detailData.branch }}</span></p>
        <p><span class='demo-split-pane'>新增代码行数：</span><span class="demo-split">{{ detailData.sumlines }}</span></p>
        <p><span class='demo-split-pane'>commitId：</span><span class="demo-split">{{ detailData.commit_id }}</span></p>
        <p><span class='demo-split-pane'>提交人：</span><span class="demo-split">{{ detailData.commit_user }}</span></p>
        <p><span class='demo-split-pane'>提交时间：</span><span class="demo-split">{{ detailData.commit_time }}</span></p>
        <p><span class='demo-split-pane'>提交记录：</span><span class="demo-split demo-split-changelog">{{ detailData.change_log }}</span></p>
      </Card>
    </Col>
    <Col v-if="detailData.env === 'inte' || detailData.env === 'parallel'"  span="24"  style="margin-bottom: 10px">
      <card>
        <p slot="title">单元测试结果</p>
        <div class="overview-domain-measures">
            <div class="overview-domain-measure">
              <div class="overview-domain-measure-value">
                <span class="overview-domain-measure-value-text">{{detailData.utLines}}</span>
              </div>
              <div class="overview-domain-measure-label">
                <Icon type="md-snow" size="20" style="margin-right: 4px;"/>总代码行数
              </div>
            </div>
            <div class="overview-domain-measure">
              <div class="overview-domain-measure-value">
                <span class="overview-domain-measure-value-text">{{detailData.utCovered}}</span>
              </div>
              <div class="overview-domain-measure-label">
                <Icon type="logo-freebsd-devil" size="20" style="margin-right: 4px;"/>单元测试覆盖行数
              </div>
            </div>
            <div class="overview-domain-measure">
              <div class="overview-domain-measure-value">
                <a
                  :href="`http://qareport.itcjf.com/utcoverage/${detailData.app_name}/${detailData.app_version}/${detailData.app_name}`"
                  target="_blank"
                  class="overview-domain-measure-value-text">
                  {{detailData.utRate}}%
                </a>
                <!--<span class="overview-domain-measure-value-text">{{detailData.utRate}}%</span>-->
                <div class="overview-domain-measure-sup">
                  <span v-if="detailData.utRate>=80" class="rating rating-A">A</span>
                  <span v-if="detailData.utRate>=60 && detailData.utRate<80" class="rating rating-B">B</span>
                  <span v-if="detailData.utRate<60 && detailData.utRate>=40" class="rating rating-C">C</span>
                  <span v-if="detailData.utRate<40" class="rating rating-D">D</span>
                </div>
              </div>
              <div class="overview-domain-measure-label">
                <Icon type="ios-happy" size="20" style="margin-right: 4px;"/>覆盖率
              </div>
            </div>
        </div>
      </card>
    </Col>
    <Col v-if="detailData.env === 'inte' || detailData.env === 'parallel'" span="24"  style="margin-bottom: 10px">
      <card>
        <p slot="title">接口测试数据</p>
        <div class="overview-domain-measures">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <span class="overview-domain-measure-value-text">{{detailData.caseTotal}}</span>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="logo-octocat" size="20" style="margin-right: 4px;"/>接口测试用例数量
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <span class="overview-domain-measure-value-text">{{detailData.caseSuccess}}</span>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-recording" size="20" style="margin-right: 4px;"/>成功执行用例数量
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <span class="overview-domain-measure-value-text">{{detailData.passRateStr}}%</span>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.passRate>=80" class="rating rating-A">A</span>
                <span v-if="detailData.passRate>=60 && detailData.passRate<80" class="rating rating-B">B</span>
                <span v-if="detailData.passRate<60 && detailData.passRate>=40" class="rating rating-C">C</span>
                <span v-if="detailData.passRate<40" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="logo-snapchat" size="20" style="margin-right: 4px;"/>通过率
            </div>
          </div>
        </div>
      </card>
    </Col>
    <Col v-if="detailData.env === 'inte' || detailData.env === 'parallel'" span="24" style="margin-bottom: 10px">
      <Card>
        <p slot="title">sonar检测结果</p>
        <div class="overview-domain-measures" style="margin-bottom: 30px">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=BUG`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.bugs}}
              </a>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.bugs<1" class="rating rating-A">A</span>
                <span v-if="detailData.bugs>=1 && detailData.bugs<5" class="rating rating-B">B</span>
                <span v-if="detailData.bugs>=5 && detailData.bugs<=10" class="rating rating-C">C</span>
                <span v-if="detailData.bugs>10" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-bug" size="20" style="margin-right: 4px;"/>bugs
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=BUG`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.new_bugs}}
              </a>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.new_bugs<1" class="rating rating-A">A</span>
                <span v-if="detailData.new_bugs>=1 && detailData.new_bugs<5" class="rating rating-B">B</span>
                <span v-if="detailData.new_bugs>=5 && detailData.new_bugs<10" class="rating rating-C">C</span>
                <span v-if="detailData.new_bugs>=10" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-bug" size="20" style="margin-right: 4px;"/>new_bugs
            </div>
          </div>
        </div>
        <div class="overview-domain-measures" style="margin-bottom: 30px">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=VULNERABILITY`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.vulnerabilities}}
              </a>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.vulnerabilities<20" class="rating rating-A">A</span>
                <span v-if="detailData.vulnerabilities>=20 && detailData.vulnerabilities<70" class="rating rating-B">B</span>
                <span v-if="detailData.vulnerabilities>=70 && detailData.vulnerabilities<100" class="rating rating-C">C</span>
                <span v-if="detailData.vulnerabilities>=100" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="md-lock" size="20" style="margin-right: 4px;"/>vulnerabilities
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=VULNERABILITY`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.new_vulnerabilities}}
              </a>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.new_vulnerabilities<20" class="rating rating-A">A</span>
                <span v-if="detailData.new_vulnerabilities>=20 && detailData.new_vulnerabilities<70" class="rating rating-B">B</span>
                <span v-if="detailData.new_vulnerabilities>=70 && detailData.new_vulnerabilities<100" class="rating rating-C">C</span>
                <span v-if="detailData.new_vulnerabilities>=100" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="md-lock" size="20" style="margin-right: 4px;"/>new_vulnerabilities
            </div>
          </div>
        </div>
        <div class="overview-domain-measures" style="margin-bottom: 30px">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/component_measures?id=${detailData.app_name}%3A${detailData.branch}&metric=duplicated_blocks`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.duplicated_blocks}}
              </a>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="md-contacts" size="20" style="margin-right: 4px;"/>duplicated_blocks
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/component_measures?id=${detailData.app_name}%3A${detailData.branch}&metric=duplicated_lines_density`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.duplications}}%
              </a>
              <div class="overview-domain-measure-sup">
                <span v-if="detailData.duplications<20" class="rating rating-A">A</span>
                <span v-if="detailData.duplications>=20 && detailData.duplications<40" class="rating rating-B">B</span>
                <span v-if="detailData.duplications>=40 && detailData.duplications<60" class="rating rating-C">C</span>
                <span v-if="detailData.duplications>=60" class="rating rating-D">D</span>
              </div>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="md-contacts" size="20" style="margin-right: 4px;"/>duplications
            </div>
          </div>
        </div>
        <div class="overview-domain-measures" style="margin-bottom: 30px">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=CODE_SMELL`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.code_smells}}
              </a>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-nuclear" size="20" style="margin-right: 4px;"/>code_smells
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <a
                :href="`http://10.10.0.199:9091/project/issues?id=${detailData.app_name}%3A${detailData.branch}&resolved=false&types=CODE_SMELL`"
                target="_blank"
                class="overview-domain-measure-value-text">
                {{detailData.new_code_smells}}
              </a>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-nuclear" size="20" style="margin-right: 4px;"/>new_code_smells
            </div>
          </div>
        </div>
        <div class="overview-domain-measures" style="margin-bottom: 30px">
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <span class="overview-domain-measure-value-text">{{detailData.reliability_rating}}</span>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-paper-plane" size="20" style="margin-right: 4px;"/>reliability_rating
            </div>
          </div>
          <div class="overview-domain-measure">
            <div class="overview-domain-measure-value">
              <span class="overview-domain-measure-value-text">{{detailData.security_rating}}</span>
            </div>
            <div class="overview-domain-measure-label">
              <Icon type="ios-paper-plane" size="20" style="margin-right: 4px;"/>security_rating
            </div>
          </div>
        </div>
      </Card>
    </Col>
  </Row>

</template>

<script>
import { httpRequest } from '@/libs/tools.js'
import { getBuildInfo } from '@/api/deploy.js'
import '../css/deploy.less'
export default {
  name: 'buildDetailsNew',
  data () {
    return {
      detailData: []
    }
  },
  methods: {
    buildDetailsNew (app_version) {
      let params = {
        app_version: app_version
      }
      this.showDetail = true
      httpRequest(this, getBuildInfo, [params]).then(res => {
        console.log('res', res)
        this.detailData = res.result[0]
        // this.deployDetail = res.deploy
      })
    },
    getColor (rate) {
      let color
      if (rate >= 80) {
        color = '#0a0'
      } else if (rate >= 60 && rate < 80) {
        color = '#eabe06'
      } else if (rate < 60) {
        color = '#e00'
      }
      return color
    }
  }

}
</script>

<style scoped>
</style>
