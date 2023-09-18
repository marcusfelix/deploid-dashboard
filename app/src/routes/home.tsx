import { useQuery } from "react-query";
import SideBar from "../components/side-bar";
import { client } from "../client";
import Router, { Route } from "preact-router";
import NavBar from "../components/nav-bar";
import NotFound from "./not-found";
import Secrets from "./secrets";
import Releases from "./releases";
import Deploy from "./deploy";
import Tasks from "./tasks";
import { Record } from "pocketbase";
import Project from "./project";
import Subscription from "./subscription";

export default function Home(props: any){
    const user = client.authStore.model as Record
    const slug = props.slug.split("/").at(0);
    const action = props.slug.split("/").at(-1);
    
    const query = useQuery(`projects-${user.id}`, async () => client.collection("projects").getList());
    const project = (query.data?.items ?? []).find((e) => e.slug === slug);

    return <div class="h-screen flex gap-4 px-4 tracking-normal overflow-y-hidden snap-x">
        <SideBar projects={query.data?.items ?? []} slug={slug}/>
        {project ? <NavBar project={project} action={action}/> : null}
        {project ? <Router>
            <Route path="/:slug/" project={project} component={Project}/>
            <Route path="/:slug/tasks" project={project} component={Tasks}/>
            <Route path="/:slug/secrets" project={project} component={Secrets}/>
            <Route path="/:slug/releases" project={project} component={Releases}/>
            <Route path="/:slug/subscription" project={project} component={Subscription}/>
            <Route path="/deploy/:template" component={Deploy}/>
            <Route default component={NotFound}/>
        </Router> : <Router>
            <Route path="/deploy/:template" component={Deploy}/>
            <Route default component={NotFound}/>
        </Router>}
        {project ? null : null}
    </div>
}