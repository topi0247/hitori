import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Provider, createStore } from "jotai";
import type { ReactNode } from "react";
import type { WritableAtom } from "jotai";
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type InitialAtomValues = Iterable<readonly [WritableAtom<any, any[], any>, any]>;

const createWrapper = (initialAtomValues?: InitialAtomValues) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  const store = createStore();
  if (initialAtomValues) {
    for (const [atom, value] of initialAtomValues) {
      store.set(atom, value);
    }
  }
  return ({ children }: { children: ReactNode }) => (
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </Provider>
  );
};

export { createWrapper };
