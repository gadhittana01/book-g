# Book Go

![cod cov](https://ffe75afa0277963feee3533053f04b47322694168a02fc7f15bc4f2-apidata.googleusercontent.com/download/storage/v1/b/gadhittana01/o/book-go-http%2Fcodcov.svg?jk=AXvcXDuzrTbroGKMftkpEZ5-_5qL2oJ3YnM7tX93-c5bGnTc_tKLylchkqHc_hh2tiHewtikY-fy1ITbJxLg08my0G4g6rBvFwmvp4r8nvxRcnpBUUo6Os6tZn9tH5E5yaZlduwS4gQ-CnP-A8MgsJXZuVZfE6K0vx6eNdL9sNlCxydUeHSamITx42bCigkcxpD313ydEDGhKOWwCFPLWVP2egE-0njsgfCnKd5poTzphmSeuV8ofAwLN_Ulq3QnwCTt4-Xwmu6fXkmq-8wwHyufviCvuqVlsx1kwnMZd3w7psqZsVdnpj05PnwzboJEDHSP4dVmwHXHSHRnzFsG5eW9CiCPH1VhEZZOgPTyAS7sZe26f-hG-M62dZIjRgZp8oPakhuXKi12UGFbPuTkK1Zk9rC4up1JRFO_Fz0ODIU2H12TkVGXBMn6gIRidTENkVC9VuPF-7qINUVkbaOF9kuJY0uFdopaY_-K8tI6R32RC64IaonttSLEIA0qRlhdYRueFcxUnAclUv09VA_5VB05UCx7i0rXTwi2dO580clLIV39_iSyL_honC3KcJvff1ZSKhZlOuRGL8N4rRNtK8DJtQaacXCrxRJAMrIYdIoZGYBYo1mtU8lbYGqDM9tZ0G53ffBdF0P3_U6RjlpdWu1w22wPbZnlob20RUSnZxhn7k7OKlCCN15yNkWiWmp6tJuiBpojWLcMhHvl7R2wdXLmeoa9b7eHt4z4NIyKfmm8jVsZM100IoxLDJzVsbL4OX4Kq4rG-H03LX6wy2k7dHoI7PwupTjMc2lKyhFYcvDKPdKZNKNniTFOjdC8DgZ0MtRVDAuiVc1IBX-E_yBB1Km2oou2p2Et3QjK3c1xWGbSglbBi2SdIV3PLp2YimMY238iUeLms2Kovej1PC-jzJrtglxHOdvZeEe-wpvzE_HrJ7SJ-Wg9YVZmHkO3Zj2ez1Y9f_tN5se_Z3tqcfxSOMbqYi1bXGxbSA45bxjv0NYKhYjZcjyCHX-y3-eFlcuYfhEiKCdN2UmGQ-j2nEZAwIr3d6_VYRfH38G6p0YsFAoBhL51zfl3yTvaDy6lfPc6WYzwYxwN0qjZPB23ZCSpUvEcfbq77phfWA_SUsqmEkFOhIWPVES0nhZoIKDJB6FlnX3d&isca=1)

## Introduction
This project is a online book store Go-based application. It uses various tools like Docker, `migrate`, `mockgen`, and `golangci-lint` to ensure smooth development, testing, and deployment workflows.

## Prerequisites to contribute

If you want to contribute to this project, ensure you have the following dependencies installed:

- [Go](https://golang.org/dl/) (v1.XX or higher)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Migrate](https://github.com/golang-migrate/migrate) (for database migrations)
- [Mockgen](https://github.com/golang/mock) (for generating mocks)
- [Wire](https://github.com/google/wire) (dependency injection)
- [Golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

Install these dependencies using the appropriate installation guides.

## Commands
1. Run the app
```bash
make start
```

2. Create a new database migration
```bash
make migrateInit name="your migration name"
```

3. Run the tests with race conditions and code coverage
```bash
make test
```

4. Generate repository mock
```bash
make mockRepo
```

5. Run the linter
```bash
make checkLint
```

## Postman documentation
https://www.postman.com/gadhittana/development/collection/0arthyu/book-go