import Router from "preact-router";
import { client } from "./client";
import Auth from "./routes/auth";
import Home from "./routes/home";
import { useEffect, useState } from "preact/hooks";
import { Admin, Record } from "pocketbase";
import { QueryClient, QueryClientProvider } from "react-query";
import Activate from "./routes/activate";

const queryClient = new QueryClient()

export function App() {
  const [user, setUser] = useState<Record | Admin | null>(client.authStore.model);

  useEffect(() => {
    client.authStore.onChange(() => setUser(client.authStore.model))
  }, []);
  
  if(!user){
    return <Auth/>
  }

  if(!user.verified){
    return <Activate/>
  }

  return <QueryClientProvider client={queryClient}>
    <Router>
      <Home path="/:slug*" />
    </Router>
  </QueryClientProvider> 
}
