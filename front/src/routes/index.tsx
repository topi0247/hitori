import {
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import { getDefaultStore } from "jotai";
import { sessionAtom } from "@/stores/auth";
import { RootPage } from "@/pages/root";
import { AuthPage } from "@/pages/auth";
import { AccountPage } from "@/pages/account";

const rootRoute = createRootRoute({
  component: Outlet,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: RootPage,
});

const authRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/auth",
  beforeLoad: () => {
    const session = getDefaultStore().get(sessionAtom);
    if (session) throw redirect({ to: "/" });
  },
  component: AuthPage,
});

const accountRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/account",
  beforeLoad: () => {
    const session = getDefaultStore().get(sessionAtom);
    if (!session) throw redirect({ to: "/auth" });
  },
  component: AccountPage,
});

const routeTree = rootRoute.addChildren([indexRoute, authRoute, accountRoute]);

const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

export { router };
