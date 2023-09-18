import { Record } from "pocketbase";
import { IconLock, IconReload } from "@tabler/icons-preact";
import { useQuery } from "react-query";
import { client } from "../client";
import { useState } from "preact/hooks";

type Props = {
  project: Record
}

export default function Subscription(props: Props) {
  const [working, setWorking] = useState(false)

  const subscription = useQuery(`subscription-${props.project.id}`, async () => client.collection("subscriptions").getFullList({
    filter: `project = '${props.project.id}'`
  }));

  const portal = async (subscription: any) => {
    setWorking(true)

    try {
      const data = await client.send(`/api/management/${props.project.id}`, {
        method: "POST",
        body: JSON.stringify({
          customer: subscription.customer,
          slug: props.project.slug,
        })
      })

      if (!data.error) {
        window.location = data.url
      }
    } catch (e) {
      console.error(e)
    }
    setWorking(false)
  }

  const color = (status: string) => {
    switch (status) {
      case "canceled":
        return "text-red-500 bg-red-100"

      case "active":
        return "text-green-500 bg-green-100"

      case "trialing":
      case "trial":
        return "text-cyan-500 bg-cyan-100"

      default:
        return "text-gray-500 bg-gray-100"
    }
  }

  if(subscription.status != "success"){
    return null
  }
  
  return <div class="w-full md:flex-1 flex flex-col gap-6 shrink-0 snap-center mx-auto max-w-4xl">
    <div class="flex justify-between py-4">
      <h1 class="text-xl font-semibold">Subscription</h1>
      <span class="bg-green-100 text-green-600 uppercase font-bold px-4 py-2 rounded-lg">{subscription.data[0].name.toUpperCase()}</span>
    </div>
    <div class="flex flex-col items-start">
      <div class="flex-1 text-lg text-gray-500">You can change between subscriptions or cancel anytime. Please, download your project data before cancelling your subscription to not lose any data.</div>
    </div>
    <div class="flex items-center gap-4">
      <button class="flex items-stretch text-xl rounded-lg overflow-hidden" disabled={subscription.status != "success"} onClick={() => portal(subscription.data[0])}><span class={`${color(subscription.data[0].status)} text-sm uppercase font-bold py-3 px-4`}>{subscription.data[0].status}</span> <span class="flex items-center gap-4 px-4 bg-gray-100 text-gray-500">Manage subscription {working ? <IconReload class="spin" /> : <IconLock/>}</span></button>
    </div>
  </div>
}