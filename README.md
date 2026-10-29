# Angry Birds Star Wars (Wii U) - NEX server

Revival server for the online leaderboards of Angry Birds Star Wars on Wii U,
run by [Protarium Network](https://github.com/Protarium-Network).

> **Credit:** this server is based on the original work of
> **Joshua Peisach ([ItzSwirlz](https://github.com/ItzSwirlz))** -
> <https://github.com/ItzSwirlz/angry-birds-star-wars>. All the original
> design, protocol work and the ranking implementation are theirs. This
> repository only adds Protarium's shared-deployment changes on top. Not an
> official Pretendo project.

Built on the [Pretendo Network](https://github.com/PretendoNetwork) NEX
libraries (`nex-go`, `nex-protocols-go`, `nex-protocols-common-go`).

## Protarium changes

- Shared-account mode: NEX passwords derived from `PN_NEX_PASSWORD_SECRET`
  (matches the account server token service), optional gRPC account lookup.
- `GetRankingsByPIDList` support and extra common-secure-server protocols.
- Fixes to ranking insert / common-data upload.

## Configuration

Everything secret comes from the environment, see `.env.example`.

## License

GPL-3.0, same as the original project. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
