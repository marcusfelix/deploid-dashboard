import { IconFolder, IconLock } from "@tabler/icons-preact";
import { IconArrowRight } from "@tabler/icons-preact";
import { IconInfoCircle } from "@tabler/icons-preact";
import { IconReload } from "@tabler/icons-preact";
import { IconCheck } from "@tabler/icons-preact";
import { useState } from "preact/hooks";
import { useQuery } from "react-query";
import { client } from "../client";
import { plans, regions } from "../includes/values";
import { Record } from "pocketbase";

export default function Deploy(props: any){
    const [working, setWorking] = useState(false);
    const [error, setError] = useState(null);
    const [name, setName] = useState("");
    const [region, setRegion] = useState("");
    const [plan, setPlan] = useState("deploy");
    const [cupon, setCupon] = useState("");

    const user = client.authStore.model as Record

    const readme = useQuery(`deploy-readme-${props.matches.template}`, async () => client.send(`/api/templates/${props.matches.template}`, {
        method: "GET"
    }), { staleTime: 1000 * 60 * 60 * 24 });

    const products = useQuery(`deploy-prices`, async () => client.send(`/api/prices`, {
        method: "GET"
    }), { staleTime: 1000 * 60 * 60 * 24 });

    const checkout = async () => {
        setWorking(true)
        try {
            const data = await client.send(`/api/checkout`, {
                method: "POST",
                body: JSON.stringify({
                    name,
                    template: props.matches.template,
                    price: getPrice().id,
                    plan: getPrice().product.name.toLowerCase(),
                    region,
                    email: user.email,
                    user: user.id,
                    cupon
                })
            })
            window.location = data.url
        } catch(e: any){
            console.log(e)
            setError(e.message)
        }
        
        setWorking(false)
    }

    const getPrice = (): any => {
        return (products.data?.data ?? []).find((e: any) => e.product.name.toLowerCase() === plan)
    }

    const validate = () => name.length > 4 && region.length > 0 && readme.data != null;

    return <div class="w-full gap-6 overflow-y-auto">
        <div class="py-4">
            <h1 class="text-xl font-semibold">Create project</h1>
        </div>
        <div class="flex flex-col flex-1 gap-4 max-w-xl">
            <div class="flex items-center gap-4 px-4 h-14 text-lg border-2 rounded-lg">
                <IconFolder/> <span class="">deploidstudio / {props.matches.template}</span>
            </div>
            <input 
                defaultValue={name}
                type="name" 
                maxLength={16}
                class="px-4 py-3 text-lg border-2 rounded-lg" 
                placeholder="Project name" 
                onInput={(e: any) => setName(e.target.value)} 
            />
            <select defaultValue={region} class="px-4 py-3 text-lg border-2 rounded-lg appearance-none" onChange={(e: any) => setRegion(e.target.value)}> 
                <option value="">Select a data center</option>
                {regions.map((region: any) => <option value={region.region}>{region.city}</option>)}
            </select>
            <div class="flex flex-col items-stretch border-2 rounded-xl overflow-hidden">
                <div class="flex bg-gray-200">
                    {Object.keys(plans).map((key: string) => <div class={`flex-1 py-3 px-4 text-lg font-bold text-center cursor-pointer ${key == plan ? 'bg-white' : 'bg-transparent text-gray-400'}`} onClick={() => setPlan(key)}>
                        {plans[key].name}
                    </div>)}
                </div>
                <div class="p-4 flex flex-col gap-4">
                    <div class="flex items-center gap-2">
                        {products.status == "success" ? <span class="bg-yellow-100 text-yellow-600 text-2xl font-bold px-3 py-1 rounded-lg">{products.status === "success" ? (getPrice().unit_amount / 100) : "-"}/{getPrice().recurring.interval[0]}</span> : null}
                    </div>
                    <ul class="flex flex-col gap-1 text-gray-500">
                        {plans[plan].services.map((service: string | string[]) => <li class="flex gap-2 items-center" key={service}>
                            {Array.isArray(service) ? <details>
                                <summary class="flex items-center gap-2 cursor-pointer"><IconCheck size={18}/> {(service as string[])[0]} <IconInfoCircle size={18}/></summary>
                                <ul class="flex flex-col py-2 pl-6 gap-1 text-gray-500">
                                    {service.slice(1).map((subservice: string) => <li class="flex items-start gap-2" key={subservice}><IconCheck size={18} class="shrink-0 mt-1"/> {subservice}</li>)}
                                </ul>
                            </details> : <div class="flex items-center gap-2"><IconCheck size={18}/> {service}</div>}
                        </li>)}
                    </ul>
                </div>
            </div>
            <input 
                type="name" 
                maxLength={16}
                class="px-4 py-3 text-lg border-2 rounded-lg" 
                placeholder="Coupon code" 
                onInput={(e: any) => setCupon(e.target.value)} 
            />
            {error ? <div class="bg-red-100 text-red-600 p-4 rounded-lg">{error}</div> : null}
            <button class={`flex justify-between px-4 py-4 text-xl ${validate() ? 'bg-black text-white' : 'bg-gray-200 text-gray-400'} rounded-xl`} disabled={!validate()} onClick={checkout}>
                {working ? <IconReload class="spin" size={26}/> : <IconLock size={26}/>}
                Checkout
                <IconArrowRight size={26}/>
            </button>
            <div class="text-sm text-gray-400">* You can upgrade and cancel your plan anytime to best meet your development demand. Learn more about <a class="text-blue-700" href="https://deploid.studio/plans" target="_blank">Deploid Studio plans</a>; ¹ All apps are deployed with our free hosting offering which are subject to our fair use policy; ² Text and assets must be provided by your business in order to be white-label eligible; ³ Service to setup and maintain an VM box of your choice with the lastest version of the app + data; ⁴ Process of build, sign and upload Appstore and Play Store apps, as well as review process follow-ups on behalf of your business. Your business must have development accounts across stores; ⁵ You can submit unlimited number os tasks where they will be analyzed and completed one-by-one. We do not take multiple tasks at the same time; ⁶ Design of branding assets for a digital product such as icons, assets and logo;</div>
            <div class="h-8"></div>
        </div>
    </div> 
}