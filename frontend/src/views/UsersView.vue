<template>
<MenuComponent :userinfo=userinfo :viewname=viewname :token=token />
Search user: <input v-on:keyup.enter="get_target_user()" v-model.lazy='account' /><button @click='get_target_user()'>Search</button>
<div v-if='target_account'>Updating user: {{ target_account }}</div>
<h2>Role</h2><hr>
Administrator <input type='checkbox' @change='modify_scopes()' v-model.lazy='target_userrole'>
<h2>Scopes</h2><hr>
<table id="datatable">
<thead>
  <tr>
    <th>Scope Name</th>
    <th>In the Scope</th>
  </tr>
</thead>
<tbody>
  <tr v-for='(scope, id) in scopes' :key=id>
    <td>{{ scope.Name }}</td>
    <td>
      <input type='checkbox' @change='modify_scopes()' v-model.lazy='scope.checked'>
    </td>
  </tr>
</tbody>
</table>
</template>

<script>
import MenuComponent from '@/components/Menu.vue'

const Role = {
  NormalUser: 0,
  Administrator: 1
};

export default {
  components: {
    MenuComponent
  },
  data() {
    return {
      viewname: 'UsersView',
      account: '',
      target_account: '',
      target_userinfo: {},
      userinfo: {},
      scopes: [],
      token: '',
    }
  },
  beforeMount() {
    this.get_scopes();
  },
  methods: {
    get_scopes: function () {
      let ref = this;
      fetch('/api/getscopes', {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText);
        ref.token = response.headers.get('X-Csrf-Token');
        return response.json();
      }).then((res) => {
        if (res.UserInfo.Role != 1) throw new Error("Only Administrator available")
        ref.userinfo = res.UserInfo;

        if (res.Scopes == null) {
          ref.scopes = [];
        }
        else {
          ref.scopes = res.Scopes;
          ref.clear_scopes_check(ref.scopes);
        }
      }).catch(function(err) {
        console.log(err);
        window.location.replace('/assets/assets.html');
      });
    },
    clear_scopes_check: function (scopes) {
      scopes.forEach((scope) => { scope.checked = false; });
    },
    get_target_user: function () {
      this.clear_scopes_check(this.scopes);

      if (!this.account.length) {
        alert("Please input target user's account")
        return;
      }

      let ref = this;
      fetch('/api/getuser_by_account?account=' + this.account, {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText);
        return response.json();
      }).then((res) => {
        ref.target_account = res.Account;
        ref.target_userinfo = res;
        ref.target_userrole = ref.target_userinfo.Role == Role.Administrator;
        ref.scopes.forEach((scope) => {
          if (ref.target_userinfo.Scopes.indexOf(scope.ID) >= 0) {
            scope.checked = true;
          }
          else {
            scope.checked = false;
          }
	});
      }).catch(function(err) {
        alert("Here is no user: " + ref.account);
        console.log(err);
      });
    },
    modify_scopes: function () {
      let new_scopes = [];
      this.scopes.forEach((scope) => {
        if (scope.checked) {
          new_scopes.push(scope.ID);
        }
      });

      if ('Scopes' in this.target_userinfo) {
        this.target_userinfo.Role = this.target_userrole ? Role.Administrator : Role.NormalUser;
        this.target_userinfo.Scopes = new_scopes;
        this.update_user_scopes();
      }
    },
    update_user_scopes: function () {
      let ref = this;
      fetch('/api/updateuser_scopes', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(ref.target_userinfo)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
      }).catch(function(err) {
        console.log(err);
      });
    },
  }
};
</script>
