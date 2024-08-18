<template>
<h2>Registration</h2>
<form id='login' method='POST' action='/api/register'>
<table class="center">
  <tr><td class="right">Account:</td><td><input type="text" name="account" required></td></tr>
  <tr><td class="right">Password:</td><td><input type="password" name="passwd" required></td></tr>
</table>
<input type="hidden" id="_csrf" name="_csrf" :value='token' />
<input type="submit" value="Register">
</form>
</template>

<script>
export default {
  data() {
    return {
      token: "",
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
