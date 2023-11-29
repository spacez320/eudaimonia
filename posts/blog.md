Blog
====

This describes how this blog is constructed. Below is a description from the bottom up.

-   Go is used to parse markdown files and generate HTML for posts. Each time a
    page is requested it will perform this conversion according to the URL
    path.

-   Gin, a Go web framework, provides URL path routing, HTML templating, and
    asset management.

-   The application runs on Kubernetes, specifically the managed Kubernetes
    offering from Linode, LKE.

-   Traefik runs in Kubernetes and performs some network routing functions,
    including:

    -   Providing NAT into the Kubernetes network for external traffic.
    -   Provisioning a Linode NodeBalancer to provision an external IP.

-   AWS's Route53 service provides DNS that resolves to the NodeBalancer
    external IP.

-   Gandi.net provides an SSL certificate which is loaded into Kubernetes as a
    Secret and terminated by Traefik as part of routing to the blog web
    service.

[![Diagram of blog components](/images/blog.png)](/images/blog.png)

Terraform provisions everything except the blog itself, which is deployed
directly via a Kubernetes manifest from GitHub Actions.

See: <https://github.com/spacez320/eudaimonia>
