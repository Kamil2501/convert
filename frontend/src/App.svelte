<script lang="ts">
import { Label, Select, Fileupload, Button, Alert, VideoPlaceholder, P } from 'flowbite-svelte';
import { convert } from './lib/convert';
import { error } from './lib/error'

let prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches

$: {
    if(prefersDark) {
        document.body.classList.add('dark')
    } else {
        if(document.body.classList.contains('dark')) {
            document.body.classList.remove('dark')
        }
    }
}

let pending = false

function setPending(val: boolean) {
    pending = val
}

let promise: Promise<Blob|undefined>

let selectedOpt: string = "";

let convertOrder: string[]

$: {
    if(selectedOpt) {
        convertOrder = selectedOpt.split(".")
    }
}


let options: {value: string, name: string}[] = [
    {value: "avi.mp4", name: "avi to mp4"},
    {value: "mp4.avi", name: "mp4 to avi"}
]

function handleSubmit(ev: Event) {
    promise = convert(ev.target as HTMLFormElement, convertOrder)
}
</script>

<main class="grid grid-cols-[0.7fr] justify-center p-5">
    <h1 class="font-bold text-[3rem]">File converter</h1>
    <Label class="flex flex-col gap-y-3 my-3">
        What type of conversion do you want to do?
        <Select items={options} bind:value={selectedOpt} />
    </Label>
    {#if convertOrder}
        <div class="grid grid-cols-2 my-5">
            <form enctype="multipart/form-data" method="post" on:submit|preventDefault={handleSubmit}>
                <Label class="flex flex-col gap-y-3">
                    Upload file ({convertOrder[0]})
                    <Fileupload name="file" accept={`.${convertOrder[0]}`} required class="w-[80%]" />
                </Label>
                <Button type="submit" class="mt-5 w-[80%]" disabled={pending}>Convert</Button>
            </form>
            {#await promise}
                <VideoPlaceholder class="ml-auto w-[400px]" />
                <p class="hidden">{setPending(true)}</p>
            {:then blob}
                {#if blob}
                    {@const videoUrl = URL.createObjectURL(blob)}
                    <div class="ml-auto">
                        <video src={videoUrl} class="w-[400px] h-[200px]" controls>
                            <track kind="captions" />
                        </video>
                        <a href={videoUrl} download><Button class="w-[400px] mt-5">Download</Button></a>
                    </div>
                    <p class="hidden">{setPending(false)}</p>
                {/if}
            {/await}
        </div>
        {#if $error}
            <Alert dismissable class="mt-5 dark:bg-[rgba(255,10,10,0.12)]">
                <p class="text-red-200">{$error}</p>
            </Alert>
        {/if}
    {/if}
</main>
