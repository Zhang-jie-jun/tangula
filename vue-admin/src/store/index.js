import Vue from 'vue'
import Vuex from 'vuex'

import state from './state'
import mutations from './mutations'
import actions from './actions'
import user from './module/user'
import app from './module/app'
import perm from './module/perm'
import deploy from './module/deploy'
import 'vue-testcase-minder-editor/lib/VueTestcaseMinderEditor.css'
import VueTestcaseMinderEditor from "_vue-testcase-minder-editor@0.3.9@vue-testcase-minder-editor";
//minder-editor
Vue.use(VueTestcaseMinderEditor)

Vue.use(Vuex)

export default new Vuex.Store({
  state,
  mutations,
  actions,
  modules: {
    user,
    app,
    perm,
    deploy,
    caseEditorStore: VueTestcaseMinderEditor.caseEditorStore
  }
})
