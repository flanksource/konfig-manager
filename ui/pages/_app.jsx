import "../styles/globals.css";
import { QueryClientProvider, QueryClient } from "react-query";
import { useRouter } from "next/router";
import { useEffect } from "react";

function App({ Component, pageProps }) {
  const queryClient = new QueryClient();

  const router = useRouter();
  useEffect(() => {
    router.push(window.location.pathname, "", { shallow: true });
    // TODO: Decide whether lint exclusion is valid
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <Component {...pageProps} />
    </QueryClientProvider>
  );
}

export default App;
