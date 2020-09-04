const state = {
	userInfo: '用户信息',
}

const mutations = {
	USERINFO: (state, userInfo) => {
		state.userInfo = userInfo
	},
}

const actions = {
	getUserInfo({ commit }, userInfo) {
		return new Promise((resolve, reject) => {
			commit('USERINFO', userInfo)
			resolve()
		})
	},
}

export default {
	namespaced: true,
	state,
	mutations,
	actions
}