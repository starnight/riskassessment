<template>
<MenuComponent :userinfo=userinfo :viewname=viewname :token=token />
Scope:
<select @change='get_assets()' v-model='chosen_scope'>
  <option v-for='(scope, id) in scopes' :value='scope.ID' :key=id>{{ scope.Name }}</option>
</select>
<table id='datatable'>
<thead>
  <tr>
    <th>Big Category</th>
    <th>Small Category</th>
    <th>Name</th>
    <th>Owner</th>
    <th>C</th>
    <th>I</th>
    <th>A</th>
    <th>Actions</th>
  </tr>
</thead>
<tbody>
  <tr v-for='(asset, id) in assets' :key=id>
  <template v-if='asset.mode == "read"'>
    <td>{{ asset.BigCategory }}</td>
    <td>{{ asset.SmallCategory }}</td>
    <td>{{ asset.Name }}</td>
    <td>{{ asset.Owner }}</td>
    <td>{{ asset.Value.Confidentiality }}</td>
    <td>{{ asset.Value.Integrity }}</td>
    <td>{{ asset.Value.Availability }}</td>
    <td>
      <button @click='edit_asset(asset)'>Edit</button>
      <button @click='delete_asset(asset)'>Delete</button>
    </td>
  </template>
  <template v-else>
    <td>
      <select v-model.lazy='asset.BigCategory'>
        <option value='' selected disabled>Choose</option>
        <option v-for='(option, id) in config.BigCategory' :value='option.value' :key=id>{{ option.value }}</option>
      </select>
    </td>
    <td>
      <select v-model.lazy='asset.SmallCategory' v-if='asset.BigCategory'>
        <option value='' selected disabled>Choose</option>
        <option v-for='(option, id) in pick_bigcategory(asset.BigCategory).SmallCategory' :value='option.value' :key=id>{{ option.value }}</option>
      </select>
    </td>
    <td><input v-model.lazy='asset.Name' /></td>
    <td><input v-model.lazy='asset.Owner' /></td>
    <td>
      <select v-model.number='asset.Value.Confidentiality'>
        <option value='' selected disabled>Choose</option>
        <option v-for='(option, id) in config.Confidentiality' :value='option.value' :key=id>{{ option.text }}</option>
      </select>
    </td>
    <td>
      <select v-model.number='asset.Value.Integrity'>
        <option value='' selected disabled>Choose</option>
        <option v-for='(option, id) in config.Integrity' :value='option.value' :key=id>{{ option.text }}</option>
      </select>
    </td>
    <td>
      <select v-model.number='asset.Value.Availability'>
        <option value='' selected disabled>Choose</option>
        <option v-for='(option, id) in config.Availability' :value='option.value' :key=id>{{ option.text }}</option>
      </select>
    </td>
    <td>
      <button @click='update_asset(asset)' v-if='asset.mode == "edit"'>Update</button>
      <button @click='add_asset(asset)' v-else>Add</button>
    </td>
  </template>
  </tr>
</tbody>
</table>
</template>

<script>
import json from '@/assets/config.json'
import MenuComponent from '@/components/Menu.vue'

export default {
  components: {
    MenuComponent
  },
  data() {
    return {
      viewname: 'AssetsView',
      userinfo: {},
      scopes: [],
      chosen_scope: '',
      assets: [],
      config: json,
      token: '',
    }
  },
  beforeMount() {
    this.get_scopes();
  },
  methods: {
    pick_bigcategory: function (category) {
      return this.config.BigCategory.find((ctg) => ctg.value == category);
    },
    get_scopes: function () {
      let ref = this;
      fetch('/api/getscopesbyuser', {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText)
        ref.token = response.headers.get('X-Csrf-Token');
        return response.json();
      }).then((res) => {
        ref.userinfo = res.UserInfo;
        ref.scopes = res.Scopes;
        ref.chosen_scope = res.Scopes[0].ID;
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    },
    get_assets: function () {
      let ref = this;
      if (!ref.chosen_scope.length) {
        alert("Please choose a scope")
        return;
      }
      fetch('/api/getassets/' + ref.chosen_scope, {
        method: 'get',
      }).then((response) => {
        if (!response.ok) throw new Error(response.statusText)
        return response.json();
      }).then((res) => {
        ref.userinfo = res.UserInfo;

        if (res.Assets == null) {
          ref.assets = [];
        }
        else {
          ref.assets = res.Assets;
          ref.assets.forEach((ast) => {
            ast.mode = 'read';
          });
        }
        ref.assets.push({
          BigCategory: ref.config.BigCategory[0].value,
          SmallCategory: '',
          Name: '',
          Owner: '',
          Value: {
            Confidentiality: ref.config.Confidentiality[0].value,
            Integrity: ref.config.Integrity[0].value,
            Availability: ref.config.Availability[0].value,
          },
          mode: 'add',
        });
      }).catch(function(err) {
        console.log(err);
      });
    },
    add_asset: function (asset) {
      let ref = this;
      asset.Scope = ref.chosen_scope;
      fetch('/api/addasset', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(asset)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    },
    edit_asset: function (asset) {
      asset.mode = 'edit'
    },
    update_asset (asset) {
      let ref = this;
      fetch('/api/updateasset', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(asset)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    },
    delete_asset: function (asset) {
      let delete_id = {ID: asset.ID};
      let ref = this;
      fetch('/api/deleteasset', {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'application/json',
          'X-Csrf-Token': ref.token,
        }),
        body: JSON.stringify(delete_id)
      }).then(function(response) {
        if (!response.ok) throw new Error(response.statusText)
        ref.get_assets();
      }).catch(function(err) {
        console.log(err);
      });
    }
  }
};
</script>
