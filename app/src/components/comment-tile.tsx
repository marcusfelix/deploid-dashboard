import ReactMarkdown from "react-markdown";

type Props = {
    comment: Comment
}

export default function CommentTile(props: Props){

    const user = () => {
        if(props.comment.body.includes("－")){
            const string = props.comment.body.substring(props.comment.body.indexOf("－ ") + 2, props.comment.body.length)
            return {
                name: string.split(" ")[0],
                login: string.split(" ")[0],
                photo: ""
            }
        }
        return props.comment.user
    }

    const body = () => {
        if(props.comment.body.includes("－")){
            return props.comment.body.substring(0, props.comment.body.indexOf("－ "),)
        }
        return props.comment.body
    }

    return <div class={`flex flex-col gap-2`}>
        <div class={`flex items-center gap-2`}>
            <div class="w-8 h-8 flex items-center justify-center font-bold text-sm aspect-square bg-red-100 text-red-800 rounded-full">
                {user().photo != "" ? <img src={user().photo} class="w-8 h-8 rounded-full object-cover"/> : <span>{(user().name ?? user().login)[0]}</span>}
            </div>
            <span class="font-bold flex-1">{user().name ?? user().login}</span>
            <span class="text-sm text-gray-400">{new Date(props.comment.created).toLocaleDateString()} @ {new Date(props.comment.created).toLocaleTimeString()}</span>
        </div>
        <div class={`prose text-gray-600`}>
            <ReactMarkdown>{body()}</ReactMarkdown>
        </div>
    </div>
}