
type Props = {
    data: string;
    extention: string;
    onChange: (data: string) => void;
}

export default function ImageEditor(props: Props){

    const handle = (e: any) => {
        const file = e.target.files[0]
        const reader = new FileReader()
        reader.onload = (e: any) => {
            props.onChange(btoa(e.target.result))
        }
        reader.readAsBinaryString(file)
    }

    return <div class="flex flex-col gap-4 py-4">
        <img src={`data:image/${props.extention};base64,${props.data}`} class="w-40 h-40 aspect-square"/>
        <input type="file" onChange={handle}/>
    </div>
}