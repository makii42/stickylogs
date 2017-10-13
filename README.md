Sticky Log Tail for docker containers
=====================================

Ever had this? 

You develop stuff locally, and your local build pipeline always restarts the docker container with the same name. Whenever that happens, you have to restart the tail on the container logs to see if stuff works. 

Now, there is [`docker events`][de] next to [`docker logs`][dl], that tells you whenever a container stops or starts. 

For now, `stickylogs` looks for events for a `container-name`, and whenever a container by that name starts, the logs of it will be streamed to the console by `stickylogs`. That will obviously stop when the container ends, but also *restart* as soon as another container with that name starts again.

Installation
------------

That's really simple: 

    $ go get -u github.com/makii42/stickylogs

Future Plans
------------

Right now, events that may trigger logs to be streamed are only filtered by container name. But the filter arguments in the docker API are much more powerful than just the name. I plan to implement arguments similar to dockers own [filters][def] at a later point. 

Contributions Welcome!

[de]: https://docs.docker.com/engine/reference/commandline/events
[dl]: https://docs.docker.com/engine/reference/commandline/logs/
[def]: https://docs.docker.com/engine/reference/commandline/events/#limiting-filtering-and-formatting-the-output