<template>
<form id='login' @submit.prevent="do_login">
<table class="center">
  <tr><td class="right">Account:</td><td><input type="text" v-model.lazy="account" required /></td></tr>
  <tr><td class="right">Password:</td><td><input type="password" v-model.lazy="passwd" required /></td></tr>
</table>
<input type="hidden" id="_csrf" name="_csrf" :value='token' />
<button type="submit">Login</button>
</form>
<br><router-link to="/registration.html">Register</router-link>
</template>

<script>
export default {
  data() {
    return {
      token: "",
      account: "",
      passwd: "",
    }
  },
  beforeMount() {
    this.get_token();
  },
  methods: {
    get_token: function () {
      let ref = this;
      fetch('/api/login', {
	method: 'get',
      }).then((response) => {
	if (!response.ok) throw new Error(response.statusText)
	ref.token = response.headers.get('X-Csrf-Token');
      }).catch(function(err) {
	console.log(err);
      });
    },
    do_login: function () {
      let formData = new FormData();
      formData.append('account', this.account);
      formData.append('passwd', this.passwd);
      formData.append('_csrf', this.token);

      fetch('/api/login', {
	method: 'post',
        body: formData,
      }).then((response) => {
	if (!response.ok) throw new Error(response.statusText)
        this.$router.push('/assets.html');
      }).catch(function(err) {
	console.log(err);
      });
    },
  }
};
</script>

<style>
table.center {
  margin-left: auto;
  margin-right: auto;
}

td.right {
  text-align: right;
}
</style>
