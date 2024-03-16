export default {
  state: {
    roleChn: '',
    userName: '',
    projects: [],
    business: []
  },
  mutations: {
    setProjects (state, params) {
      state.projects = params
    },
    setRoleChn (state, params) {
      state.roleChn = params
    },
    setPermUserName (state, params) {
      state.userName = params
    },
    setBusiness (state, params) {
      state.business = params
    }
  }
}
