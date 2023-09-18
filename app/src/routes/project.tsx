import { Record } from "pocketbase";
import { useQuery } from "react-query";
import { client } from "../client";
import BranchTile from "../components/branch-tile";

type Props = {
    project: Record
}

export default function Project(props: Props){

    const branches = useQuery(`deploy-branches-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/branches`, {
        method: "GET"
    }), {
        refetchInterval: (data): number => (data?.length ?? []) == 0 ? 5000 : 1000 * 60 * 60 * 24
    })

    const commits = useQuery(`deploy-commits-${props.project.id}`, async () => client.send(`/api/projects/${props.project.slug}/commits`, {
        method: "GET"
    }))

    return <div class="w-full md:flex-1 flex flex-col gap-6 shrink-0 snap-center mx-auto max-w-4xl">
        <div class="py-4">
            <h1 class="text-xl font-semibold">Project</h1>
        </div>
        <div class="flex flex-col gap-4">
            {(branches.data ?? []).map((branch: Branch) => <BranchTile project={props.project} branch={branch} commits={(commits.data ?? []).filter((e: any) => e.sha === branch.sha)}/>)}
        </div>
    </div>
}