<template>
	<el-form size="large" class="login-content-form">
		<el-form-item class="login-animation1">
			<el-input text placeholder="用户名" v-model="state.ruleForm.userName" clearable autocomplete="off">
				<template #prefix>
					<el-icon class="el-input__icon"><ele-User /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item class="login-animation2">
			<el-input :type="state.isShowPassword ? 'text' : 'password'" placeholder="密码"
				v-model="state.ruleForm.password" autocomplete="off">
				<template #prefix>
					<el-icon class="el-input__icon"><ele-Unlock /></el-icon>
				</template>
				<template #suffix>
					<i class="iconfont el-input__icon login-content-password"
						:class="state.isShowPassword ? 'icon-yincangmima' : 'icon-xianshimima'"
						@click="state.isShowPassword = !state.isShowPassword">
					</i>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item class="login-animation3">
			<el-col :span="15">
				<el-input text maxlength="4" :placeholder="$t('message.account.accountPlaceholder3')"
					v-model="state.ruleForm.code" clearable autocomplete="off">
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Position /></el-icon>
					</template>
				</el-input>
			</el-col>
			<el-col :span="1"></el-col>
			<el-col :span="8">
				<el-button class="login-content-code" v-waves>1234</el-button>
			</el-col>
		</el-form-item>
		<el-form-item class="login-animation4">
			<el-button type="primary" class="login-content-submit" round v-waves @click="onSignIn"
				:loading="state.loading.signIn">
				<span>登录</span>
			</el-button>
		</el-form-item>
	</el-form>
</template>

<script setup lang="ts" name="loginAccount">
	import { reactive, computed } from 'vue';
	import { useRoute, useRouter } from 'vue-router';
	import { ElMessage } from 'element-plus';
	import { useI18n } from 'vue-i18n';
	import Cookies from 'js-cookie';
	import { storeToRefs } from 'pinia';
	import { useThemeConfig } from '/@/stores/themeConfig';
	import { initFrontEndControlRoutes } from '/@/router/frontEnd';
	import { initBackEndControlRoutes } from '/@/router/backEnd';
	import { Session } from '/@/utils/storage';
	import { formatAxis } from '/@/utils/formatTime';
	import { NextLoading } from '/@/utils/loading';
	// 定义变量内容
	const { t } = useI18n();
	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const route = useRoute();
	const router = useRouter();
	const state = reactive({
		isShowPassword: false,
		ruleForm: {
			userName: 'admin',
			password: '123456',
			code: '1234',
		},
		loading: {
			signIn: false,
		},
	});

	// 登录
	const onSignIn = async () => {
		state.loading.signIn = true;
		Session.set('token', Math.random().toString(36).substr(0));
		Cookies.set('userName', state.ruleForm.userName);
		if (route.query?.redirect) {
			router.push({
				path: route.query?.redirect,
				query: Object.keys(route.query?.params).length > 0 ? JSON.parse(route.query?.params) : '',
			});
		} else {
			router.push('/');
		}
		// TODO(2023-03-16 11:24:56 谭人杰): 删除
		// if (!themeConfig.value.isRequestRoutes) {
		// 	// 前端控制路由，2、请注意执行顺序
		// 	const isNoPower = await initFrontEndControlRoutes();
		// 	signInSuccess(isNoPower);
		// } else {
		// 	// 模拟后端控制路由，isRequestRoutes 为 true，则开启后端控制路由
		// 	// 添加完动态路由，再进行 router 跳转，否则可能报错 No match found for location with path "/"
		// 	const isNoPower = await initBackEndControlRoutes();
		// 	// 执行完 initBackEndControlRoutes，再执行 signInSuccess
		// 	signInSuccess(isNoPower);
		// }
	};
	// TODO(2023-03-16 11:24:35 谭人杰): 删除
	// 登录成功后的跳转
	// const signInSuccess = (isNoPower: boolean | undefined) => {
	// 	console.log(isNoPower)
	// 	if (isNoPower) {
	// 		ElMessage.warning('抱歉，您没有登录权限');
	// 		Session.clear();
	// 	} else {
	// 		// 登录成功，跳到转首页
	// 		if (route.query?.redirect) {
	// 			router.push({
	// 				path: route.query?.redirect,
	// 				query: Object.keys(route.query?.params).length > 0 ? JSON.parse(route.query?.params) : '',
	// 			});
	// 		} else {
	// 			router.push('/');
	// 		}
	// 		// 登录成功提示
	// 		ElMessage.success(`登录成功`);
	// 		// 添加 loading，防止第一次进入界面时出现短暂空白	
	// 		NextLoading.start();
	// 	}
	// 	state.loading.signIn = false;
	// };
</script>

<style scoped lang="scss">
	.login-content-form {
		margin-top: 20px;

		.login-content-password {
			display: inline-block;
			width: 20px;
			cursor: pointer;

			&:hover {
				color: #909399;
			}
		}

		.login-content-code {
			width: 100%;
			padding: 0;
			font-weight: bold;
			letter-spacing: 5px;
		}

		.login-content-submit {
			width: 100%;
			letter-spacing: 2px;
			font-weight: 300;
			margin-top: 15px;
		}
	}
</style>