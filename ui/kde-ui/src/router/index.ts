import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: "/",
            name: "home",
            component: HomeView,
        },
        {
            path: "/about",
            name: "about",
            // route level code-splitting
            // this generates a separate chunk (About.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: () => import("../views/AboutView.vue"),
        },
        {
            path: "/dev",
            name: "dev",
            component: () => import("../views/DevSpaceStarter.vue"),
        },
        {
            path: "/devspace/:namespace/:name",
            name: "devspace",
            component: () => import("../views/DevSpaceEditor.vue"),
        },
        {
            path: "/dashboard",
            name: "dashboard",
            component: () => import("../views/DevSpaceDashboard.vue"),
        },
        {
            path: "/system",
            name: "system",
            component: () => import("../views/Installation.vue"),
        },
    ],
});

export default router;
