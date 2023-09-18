import ReactMarkdown from "react-markdown";

type Props = {
    data: string
}

export default function MarkdownViewer(props: Props){

    return <div class="prose py-10">
        <ReactMarkdown>{atob(props.data)}</ReactMarkdown>
    </div>
}