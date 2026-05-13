<template>
	<el-form size="large" class="login-content-form">
		<el-form-item class="login-animation1">
			<el-input text placeholder="请输入手机号" v-model="state.ruleForm.userName" clearable autocomplete="off">
				<template #prefix>
					<i class="iconfont icon-dianhua el-input__icon"></i>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item class="login-animation2">
			<el-col :span="15">
				<el-input text maxlength="4" placeholder="请输入手机码" v-model="state.ruleForm.code" clearable autocomplete="off">
					<template #prefix>
						<el-icon class="el-input__icon"><ele-Position /></el-icon>
					</template>
				</el-input>
			</el-col>
			<el-col :span="1"></el-col>
			<el-col :span="8">
				<el-button v-waves class="login-content-code">登录</el-button>
			</el-col>
		</el-form-item>
		<el-form-item class="login-animation3">
			<el-button round type="primary" v-waves class="login-content-submit">
				<span>登录</span>
			</el-button>
		</el-form-item>
	</el-form>
</template>

<script setup lang="ts" name="loginMobile">
import { reactive } from 'vue';

// 定义变量内容
const state = reactive({
	ruleForm: {
		userName: '',
		code: '',
	},
});
</script>

<style scoped lang="scss">
.login-content-form {
	margin-top: 20px;
	.login-content-code {
		width: 100%;
		padding: 0;
	}
	.login-content-submit {
		width: 100%;
		letter-spacing: 2px;
		font-weight: 300;
		margin-top: 15px;
	}
	.login-msg {
		color: var(--el-text-color-placeholder);
	}
}
</style>
