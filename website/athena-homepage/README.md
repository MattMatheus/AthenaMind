# Athena Homepage

Minimal static landing page scaffold for:

- `athena.teamorchestrator.com`

## Structure

- `index.html` - top-level homepage markup
- `styles/main.css` - base visual style
- `CNAME` - custom domain mapping for static hosts that support it
- `docs/` - generated at build time from repository markdown (`tools/build_docs_site.sh`)
- `robots.txt` - crawler policy
- `sitemap.xml` - initial URL map for root and docs entrypoint

## Local Preview

From this folder:

```bash
python3 -m http.server 8080
```

Then open:

- `http://127.0.0.1:8080`
