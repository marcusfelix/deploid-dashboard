import { Record } from "pocketbase";
// @ts-ignore
import Helmet from 'preact-helmet'

type Props = {
    project: Record
}

export default function Preview(props: Props){

    return <div class="flex flex-col h-screen w-full md:w-full md:max-w-sm  pb-4 shrink-0 snap-center">
        <Helmet
            script={[
                { src: `https://${props.project.slug}.fly.dev/flutter.js`, type: "text/javascript" },
                { innerHTML: `
                    window.addEventListener("load", function (ev) {
                        // Embed flutter into div#flutter_target
                        _flutter.loader.loadEntrypoint({
                        entrypointUrl: "https://${props.project.slug}.fly.dev/main.dart.js",
                        onEntrypointLoaded: async function (engine) {
                            let appRunner = await engine.initializeEngine({
                            assetBase: "https://${props.project.slug}.fly.dev/",
                            hostElement: document.querySelector("#canvas"),
                            });
                            await appRunner.runApp();
                        },
                        });
                    });
                `, type: "text/javascript" }
            ]}
        />
        <div class="flex gap-4 items-center py-4">
            <h1 class="text-2xl font-bold">Preview</h1>
            {/* <span class="bg-cyan-100 text-cyan-600 font-bold text-sm px-2 py-1 rounded-lg">DEV</span> */}
        </div>
        <div id="canvas" class="w-full h-full bg-gray-100 rounded-xl overflow-hidden"></div>
    </div>
}