<script>
  import { Router, Route } from 'svelte-routing'
  import IngestCex from './routes/IngestCex.svelte'
  import IngestImage from './routes/IngestImage.svelte'
  import PassageOverview from './routes/PassageOverview.svelte'
  import NotFound from './routes/NotFound.svelte'

  export let url = location.pathname
  console.log(url)
</script>

<svelte:head>
  <style>
    :root {
      --toolbar-bg-color: rgb(240, 240, 240);
      --toolbar-text-color: rgb(50, 50, 50);
      --toolbar-border-color: rgb(200, 200, 200);
      --pane-bg-color: rgba(248, 248, 248);
    }
  </style>
</svelte:head>

<Router {url}>
  <Route path="/tools/section/PassageOverview" component={PassageOverview} />
  <Route path="/tools/section/IngestCEX" component={IngestCex} />
  <Route path="/ingest/image" component={IngestImage} />
  <Route path="*" component={NotFound} />
</Router>
