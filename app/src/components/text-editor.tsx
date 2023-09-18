import { useEffect, useRef } from "preact/hooks";

type Props = {
    value: string;
    onChange: (data: string) => void;
}

export default function TextEditor(props: Props){
    const ref = useRef<HTMLTextAreaElement>(null)

    useEffect(() => {
        height()
    }, [props.value])

    const height = () => {
        if(ref.current){
            ref.current.style.height = 'inherit';
            ref.current.style.height = `${ref.current.scrollHeight}px`;
        }
    }

    return <textarea defaultValue={atob(props.value)} class="w-full py-4 outline-none font-mono" rows={1} onInput={(e: any) => props.onChange(btoa(e.target.value))} ref={ref}>
            
    </textarea>
}