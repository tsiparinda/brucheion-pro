<script>
  import { validateUrn } from '../lib/cts-urn'
  import PassageDesk from '../components/PassageDesk.svelte'
  import NavigationFix from '../components/NavigationFix.svelte'

  export let urn
  let passage, user, err

  $: if (!validateUrn(urn, { nid: 'cts' })) {
    urn = 'urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1'
    //  err = new Error('Passage not found')
  }

  $: Promise.all([getPassage(urn), getUser()])
    .then(([p, u]) => {
      passage = p
      user = u
    })
    .catch((e) => (err = e))

  async function getPassage(urn) {
    // const res = await fetch(`/api/v1/passage/${urn}`)
    const res = await fetch(`/api/v1/passage/undefined`)
    //const res = await fetch(`/api/v1/passage/?urn=urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.1`)
    const d = await res.json()
    return d.data
  }

  async function getUser() {
    const res = await fetch(`/api/v1/user`)
    const d = await res.json()
    return d.data
  }
  // next line was below after line <PassageDesk {passage} />     dont know why
  // <NavigationFix passageURN={passage.id} userName={user.name} />
</script>

{#if passage && !err}
  <PassageDesk {passage} />
{:else if err}
  <p>An error occurred: {err}</p>
{/if}
