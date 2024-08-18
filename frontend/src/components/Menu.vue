<template>
<ul class="menu">
  <li><router-link :class="{ active: viewname == 'AssetsView' }" to="/assets.html">Assets</router-link></li>
  <li><router-link :class="{ active: viewname == 'RiskAssessmentView' }" to="/riskassessment.html">Risk Assessment</router-link></li>
  <li><router-link :class="{ active: viewname == 'ScopesView' }" to="/scopes.html" v-if='userinfo.Role == 1'>Manage Scopes</router-link></li>
  <li><router-link :class="{ active: viewname == 'UsersView' }" to="/users.html" v-if='userinfo.Role == 1'>Manage Users</router-link></li>
  <li><a href="/api/logout">Logout</a></li>
</ul>
</template>

<script>
export default {
  name: 'MenuComponent',
  props: {
    viewname: String,
    userinfo: {},
  },
  afterMount() {
    console.log(this.props);
  },
};
</script>

<style scoped>
.menu  {
  list-style-type: none;
  margin: 0;
  padding: 0;
  overflow: hidden;
  background-color: #333;
}

.menu li {
   float: left;
}

.menu li a {
  display: block;
  color: white;
  text-align: center;
  padding: 14px 16px;
  text-decoration: none;
}

.menu li a:hover:not(.active) {
  background-color: #111;
}

.menu li a.active {
  float: right;
}

.active {
  background-color: #4CAF50;
}
</style>
