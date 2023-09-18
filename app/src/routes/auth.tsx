import { useState } from "preact/hooks"
import { client } from "../client";
import { IconAlertHexagon, IconReload } from "@tabler/icons-preact";

const regex = RegExp(/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/)

export default function Auth() {
  const [working, setWorking] = useState(false);
  const [mode, setMode] = useState<"login" | "signup">("login")
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);


  const create = async () => {
    setWorking(true);
    try {
      await client.collection("users").create({
        email,
        password,
        name,
        passwordConfirm: password
      });
      await client.collection("users").requestVerification(email)
    } catch(e: any){
      setError(e.message)
    }
    login();
    setWorking(false);
  }

  const login = async () => {
    setWorking(true);
    try {
      await client.collection("users").authWithPassword(email, password);
    } catch(e: any){
      setError(e.message)
    }
    setWorking(false);
  }

  const auth = async () => {
    setWorking(true);
    await client.collection('users').authWithOAuth2({ provider: 'google' })
    setWorking(false);
  }

  const reset = async () => {
    setWorking(true);
    await client.collection("users").requestPasswordReset(email);
    window.alert(`An e-mail was sent to ${email} with instructions to reset your password`)
    setWorking(false);
  }

  const validate = (): boolean => {
    if (mode == "login") {
      return regex.test(email) && password.length > 8;
    } else {
      return regex.test(email) && password.length > 8 && name.length > 2
    }
  }

  return <div class="w-screen h-screen flex items-center justify-center overflow-y-auto">
    <div class="w-full max-w-lg bg-white px-6 py-10 overflow-x-auto">
      {mode == "login" ? <div class="flex flex-col gap-4 justify-center">
        <img src="/logo.svg" class="w-10 h-10" />
        <h1 class="font-bold text-4xl">Welcome</h1>
        <p class="text-xl text-gray-400">Login in to access your project dashboard</p>
        <div class="h-6"></div>
        <input 
          type="email" 
          class="px-4 py-3 text-lg bg-gray-100 rounded-lg" 
          placeholder="E-mail" 
          onInput={(e: any) => setEmail(e.target.value)} 
        />
        <input 
          type="password" 
          class="px-4 py-3 text-lg bg-gray-100 rounded-lg" 
          placeholder="Password" 
          onInput={(e: any) => setPassword(e.target.value)} 
        />
        <div class="h-6"></div>
        {error ? <p class="flex gap-4 text-xl text-red-500"><IconAlertHexagon />  {error}</p> : null}
        <button 
          class={`h-14 px-4 text-xl ${validate() ? 'bg-black' : 'bg-zinc-300'}  text-white rounded-lg cursor-pointer`} 
          disabled={!validate()} 
          onClick={() => login()}
        >{working ? <IconReload class="spin mx-auto" /> : "Login"}</button>
        <button 
          class="h-14 flex gap-4 items-center justify-center p-4 text-xl rounded-lg bg-gray-100" 
          onClick={() => auth()}
        ><img src="/google.png" class="w-7 h-7"/> Continue with Google</button>
        <button 
          class="h-14 p-4 text-xl rounded-lg" 
          onClick={() => setMode("signup")}
        >Create an account</button>
        {regex.test(email) ? <button 
          class="p-4 text-xl rounded-lg" 
          onClick={() => reset()}
        >Reset password</button> : null}
      </div> : <div class="flex flex-col gap-4 justify-center">
        <img src="/logo.svg" class="w-10 h-10" />
        <h1 class="font-bold text-4xl">Create account</h1>
        <p class="text-xl text-gray-400">Provide an e-mail and password to create your account</p>
        <div class="h-6"></div>
        <input 
          type="name" 
          class="px-4 py-3 text-lg bg-gray-100 rounded-lg" 
          placeholder="Name" 
          onInput={(e: any) => setName(e.target.value)} 
        />
        <input 
          type="email" 
          class="px-4 py-3 text-lg bg-gray-100 rounded-lg" 
          placeholder="E-mail" 
          onInput={(e: any) => setEmail(e.target.value)} 
        />
        <input 
          type="password" 
          class="px-4 py-3 text-lg bg-gray-100 rounded-lg" 
          placeholder="Password" 
          onInput={(e: any) => setPassword(e.target.value)} 
        />
        <div class="h-6"></div>
        {error ? <p class="flex gap-4 text-xl text-red-500"><IconAlertHexagon />  {error}</p> : null}
        <button 
          class={`h-14 px-4 text-xl ${validate() ? 'bg-black' : 'bg-zinc-300'}  text-white rounded-lg cursor-pointer`} 
          disabled={!validate()} 
          onClick={() => create()}
        >{working ? <IconReload class="spin mx-auto" /> : "Create"}</button>
        <button 
          class="h-14 p-4 text-xl rounded-lg" 
          onClick={() => setMode("login")}
        >Login</button>
      </div>}

      <div class="flex flex-col gap-4">
        <div class="flex mt-10 gap-4 opacity-30 items-start justify-between">
          <div class="">
            <a href="https://deploid.studio" class="flex items-center gap-3">
              <img src="/logo.svg" class="h-7 w-7"/>
              <p class="text-xl font-bold tracking-tight">Deplōid Studio</p>
            </a>
          </div>
        </div>
        <div class="text-sm opacity-30">
          © {new Date().getFullYear()} Deplōid Studio • <a href="/terms">Terms of service</a> • <a href="/privacy">Privacy policy</a>
        </div>
      </div>
    </div>
  </div>
}