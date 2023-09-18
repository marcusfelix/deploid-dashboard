import { IconMoodSmile } from "@tabler/icons-preact";
import { Record } from "pocketbase";
import useHover from "../includes/use-hover";
import { client } from "../client";
import { useQueryClient } from "react-query";

type Props = {
    message: Record
}

export default function Reactions(props: Props){
    const query = useQueryClient()
    const [hoverRef, isHovered]: any = useHover();

    const user = client.authStore.model as Record

    const reactions = {
        
        "👍": "thumbs up",
        "👎": "thumbs down",
        "😎": "cool",
        "👀": "eyes",
        "🚀": "rocket",
        "🎉": "party",
        "💩": "poop",
    }

    const add = async (reaction: string) => {
        await client.collection("messages").update(props.message.id, {
            reactions: Object.assign({}, props.message.reactions ?? {}, {
                [reaction]: Object.assign([], (props.message.reactions?.[reaction] ?? []), [user.id])
            })
        })

        // Invalidade messages
        query.invalidateQueries({ queryKey: [`messages-${props.message.project}`] })
    }

    return <div class="flex px-1 pb-4 gap-3 items-center">
        {Object.keys(props.message.reactions ?? {}).map((reaction) => <span class="">
            <span class="relative cursor-pointer grayscale hover:grayscale-0" onClick={() => add(reaction)} >
                <span class="text-sm">{reaction}</span>
                {props.message.reactions[reaction].length > 1 ? <span class="px-[6px] min-w-4 h-4 flex items-center justify-center bg-gray-300 text-gray-600 absolute -top-2 -right-2 text-[11px] font-bold rounded-full">{props.message.reactions[reaction].length}</span> : null}
            </span>
        </span>)}
        <div class="flex gap-3">
            <div class="flex gap-3 relative" ref={hoverRef}>
                <IconMoodSmile class="text-gray-400 hover:text-black cursor-pointer" size={22}/>
                <div class={`absolute -left-2 -bottom-3 bg-black flex flex-col gap-2 px-2 py-3 rounded-full ${isHovered ? 'opacity-100' : 'opacity-0 h-0 overflow-hidden'}`}>
                    {Object.keys(reactions).map((reaction) => <span onClick={() => add(reaction)} class="text grayscale hover:grayscale-0 cursor-pointer">{reaction}</span>)}
                    <IconMoodSmile class="text-white" size={22}/>
                </div>
            </div>
        </div>
        
    </div>
}