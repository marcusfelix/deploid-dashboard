import { IconCircleDot, IconExternalLink, IconReload } from "@tabler/icons-preact"
import { Record } from "pocketbase";
import { useQuery } from "react-query";
import { client } from "../client";

type Props = {
  project: Record,
  branch: Branch,
  commits: Commit[]
}


export default function BranchTile(props: Props){

  const refetch = (data: any): number => (data ?? []).find((e: any) => e.status == "in_progress") != null ? 10000 : (1000 * 60 * 60 * 24)

  const query = useQuery(`workflows-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/workflow/${props.branch.name}`, {
      method: "GET"
  }), {
    refetchInterval: (data) => refetch(data)
  });

  const translate = (key: string) => {
    switch(key){
      case "main": return "Production"
      case "dev": return "Development"
      default: return key
    }
  }

  const msToTime = (duration: number) => {
    var seconds = Math.floor((duration / 1000) % 60),
      minutes = Math.floor((duration / (1000 * 60)) % 60);
  
    return minutes + "m " + seconds + "s";
  }

  const isMain = props.branch.name === "main"
  const isLive = query.data?.find((e: any) => e.status == "completed" && e.name === "Deploy") != null
  
  return <div class="flex flex-col bg-gray-50 border-2 border-gray-200 rounded-xl overflow-hidden">
    <div class={"p-4 flex flex-col gap-6 bg-white" + (isMain ? "" : "")}>
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-semibold">{translate(props.branch.name)}</h2>
        <div class={"px-2 py-1 font-bold uppercase text-sm rounded-lg "  + (isMain ? "bg-green-200 text-green-700" : "bg-gray-300 text-gray-700")}>{props.branch.name}</div>
      </div>
      <div class="flex flex-col gap-2">
        {props.commits.map((commit: Commit) => <div class="flex flex-col md:flex-row md:items-start gap-2 md:gap-4 font-mono text-gray-500">
          <span class="flex flex-1 items-center gap-4">
            <div class={`bg-black/20 w-6 h-6 rounded-full overflow-hidden shrink-0`} alt={commit.author.name}>
              <img src={`${commit.author.avatar}`} class="object-fit"/>
            </div>
            <span class="text-sm">{commit.message} • <span class="text-gray-500 text-sm">{new Date(commit.updated).toLocaleString()}</span></span>
          </span>
          
        </div>)}
      </div>
      <div class="flex">
        <a class={"flex items-center gap-2 px-4 py-2 text-white text-lg rounded-lg " + (isLive ? "bg-blue-600" : "bg-gray-400")} href={"https://" + props.project.slug + ".fly.dev"} target="_blank"><IconExternalLink/> Open</a>
      </div>
    </div>
    {query.status == "success" ? <div class="flex flex-col px-4 py-1">
      {(query.data ?? []).map((workflow: Workflow, i: number) => <div class="flex gap-2 items-center text-gray-400">
        <div class="w-8 flex justify-center items-stretch py-2 relative">
          <div class="bg-gray-50 rounded-full p-1 relative z-10 text-gray-300">
            {workflow.status == "in_progress" ? <IconReload class="slow-spin" size={22}/> : workflow.status == "completed" ? <IconCircleDot class="text-green-500" size={22}/> : workflow.status == "in-failure" ? <IconCircleDot class="text-red-500" size={22}/> : <IconCircleDot class="text-gray-400" size={22}/>}
          </div>
          {i < ((query.data ?? []).length - 1) ? <div class="w-[2px] h-full absolute bg-gray-300"></div> : null}
        </div>
        <div class="flex flex-col flex-1">
          <p>{workflow.name}</p>
        </div>
        {workflow.status == "in_progress" ? <span class="text-sm">{msToTime(new Date().getTime() - new Date(workflow.created).getTime())}</span> : null}
      </div>)}
    </div> : null}
  </div>
}