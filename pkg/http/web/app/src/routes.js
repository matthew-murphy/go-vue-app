import Vue from 'vue';
import VueRouter from 'vue-router';
import Auth from '@okta/okta-vue';

import Home from './components/Home';
import Login from './components/Login';
import GitHubRepoDetails from './components/GitHubRepoDetails';

Vue.use(ViewRouter);
Vue.use(Auth, {
    issurer: 'dev-24316545.okta.com/oauth2/default',
    client_id: '0oa3t3c6k7sVQwO7C5d7',
    redirect_uri: 'http://localhost:8080/implicit/callback',
    scope: 'openid profile email'
})

export default new VueRouter({
    mode: 'history',
    routes: [
        { path: '/', component: Login },
        { path: '/me', component: Home, meta: { requiresAuth: true }},
        { name: 'repo-details', path: '/repo/:id', component: GitHubRepoDetails, meta: { requiresAuth: true } },
        { path: '/implicit/callback', component: Auth.handleCallback() }
    ]
});